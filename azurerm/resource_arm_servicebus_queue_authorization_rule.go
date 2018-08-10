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

func resourceArmServiceBusQueueAuthorizationRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmServiceBusQueueAuthorizationRuleCreateUpdate,
		Read:   resourceArmServiceBusQueueAuthorizationRuleRead,
		Update: resourceArmServiceBusQueueAuthorizationRuleCreateUpdate,
		Delete: resourceArmServiceBusQueueAuthorizationRuleDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(time.Minute * 30),
			Update: schema.DefaultTimeout(time.Minute * 30),
			Delete: schema.DefaultTimeout(time.Minute * 30),
		},

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

			"queue_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateServiceBusQueueName(),
			},

			"resource_group_name": resourceGroupNameSchema(),
		}),

		CustomizeDiff: azure.ServiceBusAuthorizationRuleCustomizeDiff,
	}
}

func resourceArmServiceBusQueueAuthorizationRuleCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).serviceBusQueuesClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM ServiceBus Queue Authorization Rule creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	namespaceName := d.Get("namespace_name").(string)
	queueName := d.Get("queue_name").(string)

	if d.IsNewResource() {
		// first check if there's one in this subscription requiring import
		resp, err := client.GetAuthorizationRule(ctx, resourceGroup, namespaceName, queueName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error checking for the existence of Service Bus Queue Authorization Rule %q (Resource Group %q / Namespace %q): %+v", name, resourceGroup, namespaceName, err)
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
	if _, err := client.CreateOrUpdateAuthorizationRule(waitCtx, resourceGroup, namespaceName, queueName, name, parameters); err != nil {
		return fmt.Errorf("Error creating/updating ServiceBus Queue Authorization Rule %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.GetAuthorizationRule(ctx, resourceGroup, namespaceName, queueName, name)
	if err != nil {
		return err
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read ServiceBus Namespace Queue Authorization Rule %q (Queue %q / Namespace %q / Resource Group %q) ID", name, queueName, namespaceName, resourceGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmServiceBusQueueAuthorizationRuleRead(d, meta)
}

func resourceArmServiceBusQueueAuthorizationRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).serviceBusQueuesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	namespaceName := id.Path["namespaces"]
	name := id.Path["authorizationRules"]
	queueName := id.Path["queues"]

	resp, err := client.GetAuthorizationRule(ctx, resGroup, namespaceName, queueName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure ServiceBus Queue Authorization Rule %q (Queue %q / Namespace %q / Resource Group %q): %+v", name, queueName, namespaceName, resGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("namespace_name", namespaceName)
	d.Set("queue_name", queueName)

	if properties := resp.SBAuthorizationRuleProperties; properties != nil {
		listen, send, manage := azure.FlattenServiceBusAuthorizationRuleRights(properties.Rights)
		d.Set("manage", manage)
		d.Set("listen", listen)
		d.Set("send", send)
	}

	keysResp, err := client.ListKeys(ctx, resGroup, namespaceName, queueName, name)
	if err != nil {
		return fmt.Errorf("Error making Read request on Azure ServiceBus Queue Authorization Rule List Keys %q: %+v", name, err)
	}

	d.Set("primary_key", keysResp.PrimaryKey)
	d.Set("primary_connection_string", keysResp.PrimaryConnectionString)
	d.Set("secondary_key", keysResp.SecondaryKey)
	d.Set("secondary_connection_string", keysResp.SecondaryConnectionString)

	return nil
}

func resourceArmServiceBusQueueAuthorizationRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).serviceBusQueuesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	namespaceName := id.Path["namespaces"]
	name := id.Path["authorizationRules"]
	queueName := id.Path["queues"]

	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutDelete))
	defer cancel()
	if _, err = client.DeleteAuthorizationRule(waitCtx, resGroup, namespaceName, queueName, name); err != nil {
		return fmt.Errorf("Error issuing delete request of ServiceBus Queue Authorization Rule %q (Queue %q / Namespace %q / Resource Group %q): %+v", name, queueName, namespaceName, resGroup, err)
	}

	return nil
}
