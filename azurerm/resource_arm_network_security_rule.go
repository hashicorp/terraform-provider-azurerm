package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/arm/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmNetworkSecurityRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmNetworkSecurityRuleCreate,
		Read:   resourceArmNetworkSecurityRuleRead,
		Update: resourceArmNetworkSecurityRuleCreate,
		Delete: resourceArmNetworkSecurityRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"network_security_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLength(140),
			},

			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.SecurityRuleProtocolAsterisk),
					string(network.SecurityRuleProtocolTCP),
					string(network.SecurityRuleProtocolUDP),
				}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},

			"source_port_range": {
				Type:     schema.TypeString,
				Required: true,
			},

			"destination_port_range": {
				Type:     schema.TypeString,
				Required: true,
			},

			"source_address_prefix": {
				Type:     schema.TypeString,
				Required: true,
			},

			"destination_address_prefix": {
				Type:     schema.TypeString,
				Required: true,
			},

			"access": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.SecurityRuleAccessAllow),
					string(network.SecurityRuleAccessDeny),
				}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},

			"priority": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(100, 4096),
			},

			"direction": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.SecurityRuleDirectionInbound),
					string(network.SecurityRuleDirectionOutbound),
				}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},
		},
	}
}

func resourceArmNetworkSecurityRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).secRuleClient

	name := d.Get("name").(string)
	nsgName := d.Get("network_security_group_name").(string)
	resGroup := d.Get("resource_group_name").(string)

	source_port_range := d.Get("source_port_range").(string)
	destination_port_range := d.Get("destination_port_range").(string)
	source_address_prefix := d.Get("source_address_prefix").(string)
	destination_address_prefix := d.Get("destination_address_prefix").(string)
	priority := int32(d.Get("priority").(int))
	access := d.Get("access").(string)
	direction := d.Get("direction").(string)
	protocol := d.Get("protocol").(string)

	azureRMLockByName(nsgName, networkSecurityGroupResourceName)
	defer azureRMUnlockByName(nsgName, networkSecurityGroupResourceName)

	rule := network.SecurityRule{
		Name: &name,
		SecurityRulePropertiesFormat: &network.SecurityRulePropertiesFormat{
			SourcePortRange:          &source_port_range,
			DestinationPortRange:     &destination_port_range,
			SourceAddressPrefix:      &source_address_prefix,
			DestinationAddressPrefix: &destination_address_prefix,
			Priority:                 &priority,
			Access:                   network.SecurityRuleAccess(access),
			Direction:                network.SecurityRuleDirection(direction),
			Protocol:                 network.SecurityRuleProtocol(protocol),
		},
	}

	if v, ok := d.GetOk("description"); ok {
		description := v.(string)
		rule.SecurityRulePropertiesFormat.Description = &description
	}

	_, createErr := client.CreateOrUpdate(resGroup, nsgName, name, rule, make(chan struct{}))
	err := <-createErr
	if err != nil {
		return err
	}

	read, err := client.Get(resGroup, nsgName, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Security Group Rule %s/%s (resource group %s) ID", nsgName, name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmNetworkSecurityRuleRead(d, meta)
}

func resourceArmNetworkSecurityRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).secRuleClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	networkSGName := id.Path["networkSecurityGroups"]
	sgRuleName := id.Path["securityRules"]

	resp, err := client.Get(resGroup, networkSGName, sgRuleName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure Network Security Rule %q: %+v", sgRuleName, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)

	if props := resp.SecurityRulePropertiesFormat; props != nil {
		d.Set("access", string(props.Access))
		d.Set("destination_address_prefix", props.DestinationAddressPrefix)
		d.Set("destination_port_range", props.DestinationPortRange)
		d.Set("direction", string(props.Direction))
		d.Set("description", props.Description)
		d.Set("priority", int(*props.Priority))
		d.Set("protocol", string(props.Protocol))
		d.Set("source_address_prefix", props.SourceAddressPrefix)
		d.Set("source_port_range", props.SourcePortRange)
	}

	return nil
}

func resourceArmNetworkSecurityRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).secRuleClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	nsgName := id.Path["networkSecurityGroups"]
	sgRuleName := id.Path["securityRules"]

	azureRMLockByName(nsgName, networkSecurityGroupResourceName)
	defer azureRMUnlockByName(nsgName, networkSecurityGroupResourceName)

	_, deleteErr := client.Delete(resGroup, nsgName, sgRuleName, make(chan struct{}))
	err = <-deleteErr

	return err
}
