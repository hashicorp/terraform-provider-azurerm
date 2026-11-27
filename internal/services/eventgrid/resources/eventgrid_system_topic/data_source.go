// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_system_topic

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/systemtopics"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	eventgridClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/eventgrid/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Event Grid System Topic.",
		Read:        dataSourceEventGridSystemTopicRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description: "The name of the Event Grid System Topic resource.",
				Type:        pluginsdk.TypeString,
				Required:    true,
				ValidateFunc: validation.All(
					validation.StringIsNotEmpty,
					validation.StringMatch(
						regexp.MustCompile("^[-a-zA-Z0-9]{3,128}$"),
						"EventGrid Topics name must be 3 - 128 characters long, contain only letters, numbers and hyphens.",
					),
				),
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"identity": commonschema.SystemOrUserAssignedIdentityComputed(),

			"location": commonschema.LocationComputed(),

			"source_resource_id": {
				Description: "The ID of the Event Grid System Topic ARM Source.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"topic_type": {
				Description: "The Topic Type of the Event Grid System Topic.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"metric_resource_id": {
				Description: "The Metric Resource ID of the Event Grid System Topic.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceEventGridSystemTopicRead(d *pluginsdk.ResourceData, meta any) error {
	client := eventgridClient.From(meta).SystemTopics
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id := systemtopics.NewSystemTopicID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.SystemTopicName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			d.Set("metric_resource_id", props.MetricResourceId)
			d.Set("source_resource_id", props.Source)
			d.Set("topic_type", props.TopicType)
		}

		flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", flattenedIdentity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	return nil
}
