// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package digitaltwins

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/digitaltwins/2023-01-31/digitaltwinsinstance"
	"github.com/hashicorp/go-azure-sdk/resource-manager/digitaltwins/2023-01-31/endpoints"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/digitaltwins/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDigitalTwinsEndpointEventGrid() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDigitalTwinsEndpointEventGridCreateUpdate,
		Read:   resourceDigitalTwinsEndpointEventGridRead,
		Update: resourceDigitalTwinsEndpointEventGridCreateUpdate,
		Delete: resourceDigitalTwinsEndpointEventGridDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := endpoints.ParseEndpointID(id)
			return err
		}, validateEndpointType(func(input endpoints.DigitalTwinsEndpointResourceProperties) error {
			if _, ok := input.(endpoints.EventGrid); !ok {
				return fmt.Errorf("expected an EventGrid type but got: %+v", input)
			}
			return nil
		})),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DigitalTwinsInstanceName,
			},

			"digital_twins_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: digitaltwinsinstance.ValidateDigitalTwinsInstanceID,
			},

			"eventgrid_topic_endpoint": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.IsURLWithHTTPS,
			},

			"eventgrid_topic_primary_access_key": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"eventgrid_topic_secondary_access_key": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"dead_letter_storage_secret": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceDigitalTwinsEndpointEventGridCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).DigitalTwins.EndpointClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	digitalTwinsId, err := digitaltwinsinstance.ParseDigitalTwinsInstanceID(d.Get("digital_twins_id").(string))
	if err != nil {
		return err
	}

	id := endpoints.NewEndpointID(subscriptionId, digitalTwinsId.ResourceGroupName, digitalTwinsId.DigitalTwinsInstanceName, d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.DigitalTwinsEndpointGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of an existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_digital_twins_endpoint_eventgrid", id.ID())
		}
	}

	payload := endpoints.DigitalTwinsEndpointResource{
		Properties: &endpoints.EventGrid{
			AuthenticationType: pointer.To(endpoints.AuthenticationTypeKeyBased),
			TopicEndpoint:      d.Get("eventgrid_topic_endpoint").(string),
			AccessKey1:         d.Get("eventgrid_topic_primary_access_key").(string),
			AccessKey2:         pointer.To(d.Get("eventgrid_topic_secondary_access_key").(string)),
			DeadLetterSecret:   pointer.To(d.Get("dead_letter_storage_secret").(string)),
		},
	}

	if err := client.DigitalTwinsEndpointCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceDigitalTwinsEndpointEventGridRead(d, meta)
}

func resourceDigitalTwinsEndpointEventGridRead(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).DigitalTwins.EndpointClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := endpoints.ParseEndpointID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.DigitalTwinsEndpointGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	d.Set("name", id.EndpointName)
	d.Set("digital_twins_id", digitaltwinsinstance.NewDigitalTwinsInstanceID(subscriptionId, id.ResourceGroupName, id.DigitalTwinsInstanceName).ID())

	if model := resp.Model; model != nil {
		props, ok := model.Properties.(endpoints.EventGrid)
		if !ok {
			return fmt.Errorf("retrieving %s: expected an EventGrid type but got: %+v", *id, model.Properties)
		}

		d.Set("eventgrid_topic_endpoint", props.TopicEndpoint)
	}

	return nil
}

func resourceDigitalTwinsEndpointEventGridDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DigitalTwins.EndpointClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := endpoints.ParseEndpointID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DigitalTwinsEndpointDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
