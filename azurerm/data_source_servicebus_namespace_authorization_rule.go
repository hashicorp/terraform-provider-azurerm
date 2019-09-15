package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmServiceBusNamespaceAuthorizationRule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmServiceBusNamespaceAuthorizationRuleRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"namespace_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateServiceBusNamespaceName(),
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"primary_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"primary_connection_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_connection_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceArmServiceBusNamespaceAuthorizationRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).ServiceBus.NamespacesClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	namespaceName := d.Get("namespace_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.GetAuthorizationRule(ctx, resourceGroup, namespaceName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("ServiceBus Namespace Authorization Rule %q was not found in Resource Group %q", name, resourceGroup)
		}

		return fmt.Errorf("Error retrieving ServiceBus Namespace Authorization Rule %q (Resource Group %q, Namespace %q): %s", name, resourceGroup, namespaceName, err)
	}

	d.SetId(*resp.ID)

	keysResp, err := client.ListKeys(ctx, resourceGroup, namespaceName, name)
	if err != nil {
		return fmt.Errorf("Error making Read request on Azure ServiceBus Namespace Authorization Rule List Keys %s: %+v", name, err)
	}

	d.Set("primary_key", keysResp.PrimaryKey)
	d.Set("primary_connection_string", keysResp.PrimaryConnectionString)
	d.Set("secondary_key", keysResp.SecondaryKey)
	d.Set("secondary_connection_string", keysResp.SecondaryConnectionString)

	return nil
}
