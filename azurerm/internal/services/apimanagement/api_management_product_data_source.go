package apimanagement

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2020-12-01/apimanagement"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/schemaz"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
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

			"resource_group_name": azure.SchemaResourceGroupName(),

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
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)
	productId := d.Get("product_id").(string)

	resp, err := client.Get(ctx, resourceGroup, serviceName, productId)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Product %q was not found in API Management Service %q / Resource Group %q", productId, serviceName, resourceGroup)
		}

		return fmt.Errorf("making Read request on Product %q (API Management Service %q / Resource Group %q): %+v", productId, serviceName, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	if props := resp.ProductContractProperties; props != nil {
		d.Set("approval_required", props.ApprovalRequired)
		d.Set("description", props.Description)
		d.Set("display_name", props.DisplayName)
		d.Set("published", props.State == apimanagement.Published)
		d.Set("subscriptions_limit", props.SubscriptionsLimit)
		d.Set("subscription_required", props.SubscriptionRequired)
		d.Set("terms", props.Terms)
	}

	return nil
}
