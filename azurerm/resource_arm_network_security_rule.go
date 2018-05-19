package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
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
				ValidateFunc: validation.StringLenBetween(0, 140),
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
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"source_port_ranges"},
			},

			"source_port_ranges": {
				Type:          schema.TypeSet,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Set:           schema.HashString,
				ConflictsWith: []string{"source_port_range"},
			},

			"destination_port_range": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"destination_port_ranges"},
			},

			"destination_port_ranges": {
				Type:          schema.TypeSet,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Set:           schema.HashString,
				ConflictsWith: []string{"destination_port_range"},
			},

			"source_address_prefix": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"source_address_prefixes"},
			},

			"source_address_prefixes": {
				Type:          schema.TypeSet,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Set:           schema.HashString,
				ConflictsWith: []string{"source_address_prefix"},
			},

			"destination_address_prefix": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"destination_address_prefixes"},
			},

			"destination_address_prefixes": {
				Type:          schema.TypeSet,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Set:           schema.HashString,
				ConflictsWith: []string{"destination_address_prefix"},
			},

			"source_application_security_group_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"destination_application_security_group_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
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
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	nsgName := d.Get("network_security_group_name").(string)
	resGroup := d.Get("resource_group_name").(string)

	sourcePortRange := d.Get("source_port_range").(string)
	destinationPortRange := d.Get("destination_port_range").(string)
	sourceAddressPrefix := d.Get("source_address_prefix").(string)
	destinationAddressPrefix := d.Get("destination_address_prefix").(string)
	priority := int32(d.Get("priority").(int))
	access := d.Get("access").(string)
	direction := d.Get("direction").(string)
	protocol := d.Get("protocol").(string)

	azureRMLockByName(nsgName, networkSecurityGroupResourceName)
	defer azureRMUnlockByName(nsgName, networkSecurityGroupResourceName)

	rule := network.SecurityRule{
		Name: &name,
		SecurityRulePropertiesFormat: &network.SecurityRulePropertiesFormat{
			SourcePortRange:          &sourcePortRange,
			DestinationPortRange:     &destinationPortRange,
			SourceAddressPrefix:      &sourceAddressPrefix,
			DestinationAddressPrefix: &destinationAddressPrefix,
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

	if r, ok := d.GetOk("source_port_ranges"); ok {
		var sourcePortRanges []string
		r := r.(*schema.Set).List()
		for _, v := range r {
			s := v.(string)
			sourcePortRanges = append(sourcePortRanges, s)
		}
		rule.SecurityRulePropertiesFormat.SourcePortRanges = &sourcePortRanges
	}

	if r, ok := d.GetOk("destination_port_ranges"); ok {
		var destinationPortRanges []string
		r := r.(*schema.Set).List()
		for _, v := range r {
			s := v.(string)
			destinationPortRanges = append(destinationPortRanges, s)
		}
		rule.SecurityRulePropertiesFormat.DestinationPortRanges = &destinationPortRanges
	}

	if r, ok := d.GetOk("source_address_prefixes"); ok {
		var sourceAddressPrefixes []string
		r := r.(*schema.Set).List()
		for _, v := range r {
			s := v.(string)
			sourceAddressPrefixes = append(sourceAddressPrefixes, s)
		}
		rule.SecurityRulePropertiesFormat.SourceAddressPrefixes = &sourceAddressPrefixes
	}

	if r, ok := d.GetOk("destination_address_prefixes"); ok {
		var destinationAddressPrefixes []string
		r := r.(*schema.Set).List()
		for _, v := range r {
			s := v.(string)
			destinationAddressPrefixes = append(destinationAddressPrefixes, s)
		}
		rule.SecurityRulePropertiesFormat.DestinationAddressPrefixes = &destinationAddressPrefixes
	}

	if r, ok := d.GetOk("source_application_security_group_ids"); ok {
		var sourceApplicationSecurityGroups []network.ApplicationSecurityGroup
		for _, v := range r.(*schema.Set).List() {
			sg := network.ApplicationSecurityGroup{
				ID: utils.String(v.(string)),
			}
			sourceApplicationSecurityGroups = append(sourceApplicationSecurityGroups, sg)
		}
		rule.SourceApplicationSecurityGroups = &sourceApplicationSecurityGroups
	}

	if r, ok := d.GetOk("destination_application_security_group_ids"); ok {
		var destinationApplicationSecurityGroups []network.ApplicationSecurityGroup
		for _, v := range r.(*schema.Set).List() {
			sg := network.ApplicationSecurityGroup{
				ID: utils.String(v.(string)),
			}
			destinationApplicationSecurityGroups = append(destinationApplicationSecurityGroups, sg)
		}
		rule.DestinationApplicationSecurityGroups = &destinationApplicationSecurityGroups
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, nsgName, name, rule)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating Network Security Rule %q (NSG %q / Resource Group %q): %+v", name, nsgName, resGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for completion of Network Security Rule %q (NSG %q / Resource Group %q): %+v", name, nsgName, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, nsgName, name)
	if err != nil {
		return fmt.Errorf("Error making Read request on Network Security Rule %q (NSG %q / Resource Group %q): %+v", name, nsgName, resGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Network Security Rule %s (NSG %q / resource group %s) ID", name, nsgName, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmNetworkSecurityRuleRead(d, meta)
}

func resourceArmNetworkSecurityRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).secRuleClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	networkSGName := id.Path["networkSecurityGroups"]
	sgRuleName := id.Path["securityRules"]

	resp, err := client.Get(ctx, resGroup, networkSGName, sgRuleName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Network Security Rule %q (NSG %q / Resource Group %q): %+v", sgRuleName, networkSGName, resGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("network_security_group_name", networkSGName)

	if props := resp.SecurityRulePropertiesFormat; props != nil {
		d.Set("description", props.Description)
		d.Set("protocol", string(props.Protocol))
		d.Set("destination_address_prefix", props.DestinationAddressPrefix)
		d.Set("destination_address_prefixes", props.DestinationAddressPrefixes)
		d.Set("destination_port_range", props.DestinationPortRange)
		d.Set("destination_port_ranges", props.DestinationPortRanges)
		d.Set("source_address_prefix", props.SourceAddressPrefix)
		d.Set("source_address_prefixes", props.SourceAddressPrefixes)
		d.Set("source_port_range", props.SourcePortRange)
		d.Set("source_port_ranges", props.SourcePortRanges)
		d.Set("access", string(props.Access))
		d.Set("priority", int(*props.Priority))
		d.Set("direction", string(props.Direction))
	}

	return nil
}

func resourceArmNetworkSecurityRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).secRuleClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	nsgName := id.Path["networkSecurityGroups"]
	sgRuleName := id.Path["securityRules"]

	azureRMLockByName(nsgName, networkSecurityGroupResourceName)
	defer azureRMUnlockByName(nsgName, networkSecurityGroupResourceName)

	future, err := client.Delete(ctx, resGroup, nsgName, sgRuleName)
	if err != nil {
		return fmt.Errorf("Error Deleting Network Security Rule %q (NSG %q / Resource Group %q): %+v", sgRuleName, nsgName, resGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for the deletion of Network Security Rule %q (NSG %q / Resource Group %q): %+v", sgRuleName, nsgName, resGroup, err)
	}

	return nil
}
