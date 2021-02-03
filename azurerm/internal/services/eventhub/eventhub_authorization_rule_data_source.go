package eventhub

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func EventHubAuthorizationRuleDataSource() *schema.Resource {
	return &schema.Resource{
		Read: EventHubAuthorizationRuleDataSourceRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: EventHubAuthorizationRuleSchemaFrom(map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: ValidateEventHubAuthorizationRuleName(),
			},

			"namespace_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: ValidateEventHubNamespaceName(),
			},

			"eventhub_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: ValidateEventHubName(),
			},

			"primary_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"primary_connection_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_connection_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),
		}),
	}
}

func EventHubAuthorizationRuleDataSourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.EventHubsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	eventHubName := d.Get("eventhub_name").(string)
	namespaceName := d.Get("namespace_name").(string)

	resp, err := client.GetAuthorizationRule(ctx, resourceGroup, namespaceName, eventHubName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: EventHub Authorization Rule %q (Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error: EventHub Authorization Rule %s: %+v", name, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", name)
	d.Set("eventhub_name", eventHubName)
	d.Set("namespace_name", namespaceName)
	d.Set("resource_group_name", resourceGroup)

	keysResp, err := client.ListKeys(ctx, resourceGroup, namespaceName, eventHubName, name)
	if err != nil {
		return fmt.Errorf("Error making Read request on Azure EventHub Authorization Rule List Keys %s: %+v", name, err)
	}

	d.Set("primary_key", keysResp.PrimaryKey)
	d.Set("secondary_key", keysResp.SecondaryKey)
	d.Set("primary_connection_string", keysResp.PrimaryConnectionString)
	d.Set("secondary_connection_string", keysResp.SecondaryConnectionString)
	d.Set("primary_connection_string_alias", keysResp.AliasPrimaryConnectionString)
	d.Set("secondary_connection_string_alias", keysResp.AliasSecondaryConnectionString)

	return nil
}
