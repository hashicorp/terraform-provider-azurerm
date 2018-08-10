package azurerm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmServiceBusNamespaceAuthorizationRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmServiceBusNamespaceAuthorizationRuleCreateUpdate,
		Read:   resourceArmServiceBusNamespaceAuthorizationRuleRead,
		Update: resourceArmServiceBusNamespaceAuthorizationRuleCreateUpdate,
		Delete: resourceArmServiceBusNamespaceAuthorizationRuleDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(time.Minute * 30),
			Update: schema.DefaultTimeout(time.Minute * 30),
			Delete: schema.DefaultTimeout(time.Minute * 30),
		},

		//function takes a schema map and adds the authorization rule properties to it
		Schema: azure.ServiceBusAuthorizationRuleSchemaFrom(map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateServiceBusAuthorizationRuleName(),
			},

			"namespace_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateServiceBusNamespaceName(),
			},

			"resource_group_name": resourceGroupNameSchema(),
		}),

		CustomizeDiff: azure.ServiceBusAuthorizationRuleCustomizeDiff,
	}
}

func resourceArmServiceBusNamespaceAuthorizationRuleCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).serviceBusNamespacesClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM ServiceBus Namespace Authorization Rule creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	namespaceName := d.Get("namespace_name").(string)

	if d.IsNewResource() {
		// first check if there's one in this subscription requiring import
		resp, err := client.GetAuthorizationRule(ctx, resourceGroup, namespaceName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error checking for the existence of Service Bus Namespace Authorization Rule %q (Resource Group %q / Namespace %q): %+v", name, resourceGroup, namespaceName, err)
			}
		}
		if resp.ID != nil {
			return tf.ImportAsExistsError("azurerm_servicebus_namespace_authorization_rule", *resp.ID)
		}
	}

	parameters := servicebus.SBAuthorizationRule{
		Name: &name,
		SBAuthorizationRuleProperties: &servicebus.SBAuthorizationRuleProperties{
			Rights: azure.ExpandServiceBusAuthorizationRuleRights(d),
		},
	}

	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(tf.TimeoutForCreateUpdate(d)))
	defer cancel()
	if _, err := client.CreateOrUpdateAuthorizationRule(waitCtx, resourceGroup, namespaceName, name, parameters); err != nil {
		return fmt.Errorf("Error creating/updating ServiceBus Namespace Authorization Rule %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.GetAuthorizationRule(ctx, resourceGroup, namespaceName, name)
	if err != nil {
		return err
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read ServiceBus Namespace Authorization Rule %s (resource group %s) ID", name, resourceGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmServiceBusNamespaceAuthorizationRuleRead(d, meta)
}

func resourceArmServiceBusNamespaceAuthorizationRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).serviceBusNamespacesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	namespaceName := id.Path["namespaces"]
	name := id.Path["AuthorizationRules"] //this is slightly different then a topic rule (Authorization vs authorization)

	resp, err := client.GetAuthorizationRule(ctx, resGroup, namespaceName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure ServiceBus Namespace Authorization Rule %s: %+v", name, err)
	}

	d.Set("name", name)
	d.Set("namespace_name", namespaceName)
	d.Set("resource_group_name", resGroup)

	if properties := resp.SBAuthorizationRuleProperties; properties != nil {
		listen, send, manage := azure.FlattenServiceBusAuthorizationRuleRights(properties.Rights)
		d.Set("manage", manage)
		d.Set("listen", listen)
		d.Set("send", send)
	}

	keysResp, err := client.ListKeys(ctx, resGroup, namespaceName, name)
	if err != nil {
		return fmt.Errorf("Error making Read request on Azure ServiceBus Namespace Authorization Rule List Keys %s: %+v", name, err)
	}

	d.Set("primary_key", keysResp.PrimaryKey)
	d.Set("primary_connection_string", keysResp.PrimaryConnectionString)
	d.Set("secondary_key", keysResp.SecondaryKey)
	d.Set("secondary_connection_string", keysResp.SecondaryConnectionString)

	return nil
}

func resourceArmServiceBusNamespaceAuthorizationRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).serviceBusNamespacesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	namespaceName := id.Path["namespaces"]
	name := id.Path["AuthorizationRules"] //this is slightly different then topic/queue (Authorization vs authorization)

	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutDelete))
	defer cancel()
	if _, err = client.DeleteAuthorizationRule(waitCtx, resGroup, namespaceName, name); err != nil {
		return fmt.Errorf("Error issuing Azure ARM delete request of ServiceBus Namespace Authorization Rule %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return nil
}
