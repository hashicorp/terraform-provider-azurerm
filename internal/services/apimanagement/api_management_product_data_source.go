package apimanagement

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2021-08-01/apimanagement"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewProductID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("product_id").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("making Read request on %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if props := resp.ProductContractProperties; props != nil {
		d.Set("approval_required", props.ApprovalRequired)
		d.Set("description", props.Description)
		d.Set("display_name", props.DisplayName)
		d.Set("published", props.State == apimanagement.ProductStatePublished)
		d.Set("subscriptions_limit", props.SubscriptionsLimit)
		d.Set("subscription_required", props.SubscriptionRequired)
		d.Set("terms", props.Terms)
	}

	return nil
}
