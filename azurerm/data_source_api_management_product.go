package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/apimanagement/mgmt/2018-06-01-preview/apimanagement"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceApiManagementProduct() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApiManagementProductRead,

		Schema: map[string]*schema.Schema{
			"product_id": azure.SchemaApiManagementProductDataSourceName(),

			"api_management_name": azure.SchemaApiManagementDataSourceName(),

			"resource_group_name": resourceGroupNameSchema(),

			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"subscription_required": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"approval_required": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"published": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"terms": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"subscriptions_limit": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}
func dataSourceApiManagementProductRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagementProductsClient
	ctx := meta.(*ArmClient).StopContext

	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("service_name").(string)
	productId := d.Get("product_id").(string)

	resp, err := client.Get(ctx, resourceGroup, serviceName, productId)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Product %q was not found in API Management Service %q / Resource Group %q", productId, serviceName, resourceGroup)
		}

		return fmt.Errorf("Error making Read request on Product %q (API Management Service %q / Resource Group %q): %+v", productId, serviceName, resourceGroup, err)
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
