// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_product

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/product"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	apimanagementClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing API Management Product.",
		Read:        dataSourceApiManagementProductRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"product_id": schemaz.SchemaApiManagementChildDataSourceName("The Identifier for the API Management Product."),

			"api_management_name": schemaz.SchemaApiManagementDataSourceName(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"display_name": {
				Description: "The Display Name for this API Management Product.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"subscription_required": {
				Description: "Whether a Subscription is required to access API's included in this Product.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"approval_required": {
				Description: "Whether subscribers need to be approved before being able to use the Product.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"published": {
				Description: "Whether this Product is published.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"description": {
				Description: "The description of this Product, which may include HTML formatting tags.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"terms": {
				Description: "Any Terms and Conditions for this Product, which must be accepted by Developers before they can begin the Subscription process.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"subscriptions_limit": {
				Description: "The number of subscriptions a user can have to this Product at the same time.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},
		},
	}
}

func dataSourceApiManagementProductRead(d *pluginsdk.ResourceData, meta any) error {
	client := apimanagementClient.From(meta).ProductsClient
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
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
