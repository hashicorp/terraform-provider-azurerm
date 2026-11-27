// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package digital_twins_instance

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/digitaltwins/2023-01-31/digitaltwinsinstance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	digitaltwinsClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/digitaltwins/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/digitaltwins/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Digital Twins instance.",
		Read:        dataSourceDigitalTwinsInstanceRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description:  "The name of this Digital Twins instance.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.DigitalTwinsInstanceName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"host_name": {
				Description: "The API endpoint to work with this Digital Twins instance.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceDigitalTwinsInstanceRead(d *pluginsdk.ResourceData, meta any) error {
	client := digitaltwinsClient.From(meta).InstanceClient
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id := digitaltwinsinstance.NewDigitalTwinsInstanceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.DigitalTwinsGet(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s does not exist", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.DigitalTwinsInstanceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		if props := model.Properties; props != nil {
			d.Set("host_name", props.HostName)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	d.SetId(id.ID())
	return nil
}
