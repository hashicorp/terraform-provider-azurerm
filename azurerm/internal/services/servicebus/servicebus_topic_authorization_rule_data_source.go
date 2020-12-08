package servicebus

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/servicebus/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmServiceBusTopicAuthorizationRule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmServiceBusTopicAuthorizationRuleRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.AuthorizationRuleName(),
			},

			"namespace_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NamespaceName,
			},

			"topic_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.TopicName(),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"listen": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"send": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"manage": {
				Type:     schema.TypeBool,
				Computed: true,
			},

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

func dataSourceArmServiceBusTopicAuthorizationRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.TopicsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	namespaceName := d.Get("namespace_name").(string)
	topicName := d.Get("topic_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.GetAuthorizationRule(ctx, resourceGroup, namespaceName, topicName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("ServiceBus Topic Authorization Rule %q (Resource Group %q / Namespace Name %q) was not found", name, resourceGroup, namespaceName)
		}
		return fmt.Errorf("Error making Read request on Azure ServiceBus Topic Authorization Rule %s: %+v", name, err)
	}

	d.Set("name", name)
	d.Set("topic_name", topicName)
	d.Set("namespace_name", namespaceName)
	d.Set("resource_group_name", resourceGroup)

	if properties := resp.SBAuthorizationRuleProperties; properties != nil {
		listen, send, manage := FlattenServiceBusAuthorizationRuleRights(properties.Rights)
		d.Set("listen", listen)
		d.Set("send", send)
		d.Set("manage", manage)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("API returned a nil/empty id for ServiceBus Topic Authorization Rule %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	d.SetId(*resp.ID)

	keysResp, err := client.ListKeys(ctx, resourceGroup, namespaceName, topicName, name)
	if err != nil {
		return fmt.Errorf("Error making Read request on Azure ServiceBus Topic Authorization Rule List Keys %s: %+v", name, err)
	}

	d.Set("primary_key", keysResp.PrimaryKey)
	d.Set("primary_connection_string", keysResp.PrimaryConnectionString)
	d.Set("secondary_key", keysResp.SecondaryKey)
	d.Set("secondary_connection_string", keysResp.SecondaryConnectionString)

	return nil
}
