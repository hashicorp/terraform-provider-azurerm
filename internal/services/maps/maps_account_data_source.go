// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package maps

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maps/2023-06-01/accounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/maps/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceMapsAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceMapsAccountRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.AccountName(),
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": commonschema.Tags(),

			"x_ms_client_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_access_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_access_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceMapsAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maps.AccountsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
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
