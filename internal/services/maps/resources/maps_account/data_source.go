// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package maps_account

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maps/2023-06-01/accounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	mapsClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/maps/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/maps/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Azure Maps Account.",
		Read:        dataSourceMapsAccountRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description:  "The name of the Maps Account.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.AccountName(),
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"sku_name": {
				Description: "The SKU of the Azure Maps Account.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"tags": commonschema.Tags(),

			"x_ms_client_id": {
				Description: "A unique identifier for the Maps Account.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"primary_access_key": {
				Description: "The primary key used to authenticate and authorize access to the Maps REST APIs.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
				Sensitive:   true,
			},

			"secondary_access_key": {
				Description: "The primary key used to authenticate and authorize access to the Maps REST APIs. The second key is given to provide seamless key regeneration.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func dataSourceMapsAccountRead(d *pluginsdk.ResourceData, meta any) error {
	client := mapsClient.From(meta).AccountsClient
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id := accounts.NewAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.AccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("sku_name", string(model.Sku.Name))
		if props := model.Properties; props != nil {
			d.Set("x_ms_client_id", props.UniqueId)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	keysResp, err := client.ListKeys(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving Access Keys for %s: %+v", id, err)
	}
	if model := keysResp.Model; model != nil {
		d.Set("primary_access_key", model.PrimaryKey)
		d.Set("secondary_access_key", model.SecondaryKey)
	}

	return nil
}
