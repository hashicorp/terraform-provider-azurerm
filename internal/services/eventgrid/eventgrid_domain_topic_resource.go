// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventgrid

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/domaintopics"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceEventGridDomainTopic() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceEventGridDomainTopicCreate,
		Read:   resourceEventGridDomainTopicRead,
		Delete: resourceEventGridDomainTopicDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := domaintopics.ParseDomainTopicID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringIsNotEmpty,
					validation.StringMatch(
						regexp.MustCompile("^[-a-zA-Z0-9]{3,128}$"),
						"EventGrid domain topic name must be 3 - 128 characters long, contain only letters, numbers and hyphens.",
					),
				),
			},

			"domain_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringIsNotEmpty,
					validation.StringMatch(
						regexp.MustCompile("^[-a-zA-Z0-9]{3,50}$"),
						"EventGrid domain name must be 3 - 50 characters long, contain only letters, numbers and hyphens.",
					),
				),
			},

			"resource_group_name": commonschema.ResourceGroupName(),
		},
	}
}

func resourceEventGridDomainTopicCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).EventGrid.DomainTopics
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := domaintopics.NewDomainTopicID(subscriptionId, d.Get("resource_group_name").(string), d.Get("domain_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_eventgrid_domain_topic", id.ID())
		}
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id); err != nil {
		return fmt.Errorf("creating/updating %s: %s", id, err)
	}

	d.SetId(id.ID())
	return resourceEventGridDomainTopicRead(d, meta)
}

func resourceEventGridDomainTopicRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).EventGrid.DomainTopics
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := domaintopics.ParseDomainTopicID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("%s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.TopicName)
	d.Set("domain_name", id.DomainName)
	d.Set("resource_group_name", id.ResourceGroupName)

	return nil
}

func resourceEventGridDomainTopicDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).EventGrid.DomainTopics
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := domaintopics.ParseDomainTopicID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
