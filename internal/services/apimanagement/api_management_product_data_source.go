// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/product"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceApiManagementProduct() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceApiManagementProductRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"product_id": schemaz.SchemaApiManagementChildDataSourceName(),

			"api_management_name": schemaz.SchemaApiManagementDataSourceName(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"display_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"subscription_required": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"approval_required": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"published": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"terms": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"subscriptions_limit": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceApiManagementProductRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ProductsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := product.NewProductID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("product_id").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("checking for presence of an existing %s: %+v", id, err)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("approval_required", pointer.From(props.ApprovalRequired))
			d.Set("description", pointer.From(props.Description))
			d.Set("display_name", props.DisplayName)
			d.Set("published", pointer.From(props.State) == product.ProductStatePublished)
			d.Set("subscriptions_limit", pointer.From(props.SubscriptionsLimit))
			d.Set("subscription_required", pointer.From(props.SubscriptionRequired))
			d.Set("terms", pointer.From(props.Terms))
		}
	}

	return nil
}
