package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmServiceBusTopicAuthorizationRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmServiceBusTopicAuthorizationRuleCreateUpdate,
		Read:   resourceArmServiceBusTopicAuthorizationRuleRead,
		Update: resourceArmServiceBusTopicAuthorizationRuleCreateUpdate,
		Delete: resourceArmServiceBusTopicAuthorizationRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"namespace_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"topic_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"listen": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"send": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"manage": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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

func resourceArmServiceBusTopicAuthorizationRuleCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).serviceBusTopicsClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO] preparing arguments for AzureRM ServiceBus Topic Authorization Rule creation.")

	name := d.Get("name").(string)
	namespaceName := d.Get("namespace_name").(string)
	topicName := d.Get("topic_name").(string)
	resGroup := d.Get("resource_group_name").(string)

	rights, err := expandServiceBusTopicAuthorizationRuleAccessRights(d)
	if err != nil {
		return err
	}

	parameters := servicebus.SBAuthorizationRule{
		Name: &name,
		SBAuthorizationRuleProperties: &servicebus.SBAuthorizationRuleProperties{
			Rights: rights,
		},
	}

	_, err = client.CreateOrUpdateAuthorizationRule(ctx, resGroup, namespaceName, topicName, name, parameters)
	if err != nil {
		return err
	}

	read, err := client.GetAuthorizationRule(ctx, resGroup, namespaceName, topicName, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read ServiceBus Topic Authorization Rule %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmServiceBusTopicAuthorizationRuleRead(d, meta)
}

func resourceArmServiceBusTopicAuthorizationRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).serviceBusTopicsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	namespaceName := id.Path["namespaces"]
	topicsName := id.Path["topics"]
	name := id.Path["authorizationRules"]

	resp, err := client.GetAuthorizationRule(ctx, resGroup, namespaceName, topicsName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure ServiceBus Topic Authorization Rule %s: %+v", name, err)
	}

	keysResp, err := client.ListKeys(ctx, resGroup, namespaceName, topicsName, name)
	if err != nil {
		return fmt.Errorf("Error making Read request on Azure ServiceBus Topic Authorization Rule List Keys %s: %+v", name, err)
	}

	d.Set("name", name)
	d.Set("topic_name", topicsName)
	d.Set("namespace_name", namespaceName)
	d.Set("resource_group_name", resGroup)

	flattenServiceBusTopicAuthorizationRuleAccessRights(d, resp)

	d.Set("primary_key", keysResp.PrimaryKey)
	d.Set("primary_connection_string", keysResp.PrimaryConnectionString)
	d.Set("secondary_key", keysResp.SecondaryKey)
	d.Set("secondary_connection_string", keysResp.SecondaryConnectionString)

	return nil
}

func resourceArmServiceBusTopicAuthorizationRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).serviceBusTopicsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	namespaceName := id.Path["namespaces"]
	topicName := id.Path["topics"]
	name := id.Path["authorizationRules"]

	_, err = client.DeleteAuthorizationRule(ctx, resGroup, namespaceName, topicName, name)

	if err != nil {
		return fmt.Errorf("Error issuing Azure ARM delete request of ServiceBus Topic Authorization Rule %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return nil
}

func expandServiceBusTopicAuthorizationRuleAccessRights(d *schema.ResourceData) (*[]servicebus.AccessRights, error) {
	canSend := d.Get("send").(bool)
	canListen := d.Get("listen").(bool)
	canManage := d.Get("manage").(bool)
	rights := []servicebus.AccessRights{}
	if canListen {
		rights = append(rights, servicebus.Listen)
	}

	if canSend {
		rights = append(rights, servicebus.Send)
	}

	if canManage {
		rights = append(rights, servicebus.Manage)
	}

	if len(rights) == 0 {
		return nil, fmt.Errorf("At least one Authorization Rule State must be enabled (e.g. Listen/Manage/Send)")
	}

	if canManage && !(canListen && canSend) {
		return nil, fmt.Errorf("In order to enable the 'Manage' Authorization Rule - both the 'Listen' and 'Send' rules must be enabled")
	}

	return &rights, nil
}

func flattenServiceBusTopicAuthorizationRuleAccessRights(d *schema.ResourceData, resp servicebus.SBAuthorizationRule) {

	var canListen = false
	var canSend = false
	var canManage = false

	for _, right := range *resp.Rights {
		switch right {
		case servicebus.Listen:
			canListen = true
		case servicebus.Send:
			canSend = true
		case servicebus.Manage:
			canManage = true
		default:
			log.Printf("[DEBUG] Unknown Authorization Rule Right '%s'", right)
		}
	}

	d.Set("listen", canListen)
	d.Set("send", canSend)
	d.Set("manage", canManage)
}
