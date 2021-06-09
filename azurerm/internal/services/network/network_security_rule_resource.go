package network

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-11-01/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceNetworkSecurityRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNetworkSecurityRuleCreateUpdate,
		Read:   resourceNetworkSecurityRuleRead,
		Update: resourceNetworkSecurityRuleCreateUpdate,
		Delete: resourceNetworkSecurityRuleDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"network_security_group_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 140),
			},

			"protocol": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.SecurityRuleProtocolAsterisk),
					string(network.SecurityRuleProtocolTCP),
					string(network.SecurityRuleProtocolUDP),
					string(network.SecurityRuleProtocolIcmp),
					string(network.SecurityRuleProtocolAh),
					string(network.SecurityRuleProtocolEsp),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"source_port_range": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ConflictsWith: []string{"source_port_ranges"},
			},

			"source_port_ranges": {
				Type:          pluginsdk.TypeSet,
				Optional:      true,
				Elem:          &pluginsdk.Schema{Type: pluginsdk.TypeString},
				Set:           pluginsdk.HashString,
				ConflictsWith: []string{"source_port_range"},
			},

			"destination_port_range": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ConflictsWith: []string{"destination_port_ranges"},
			},

			"destination_port_ranges": {
				Type:          pluginsdk.TypeSet,
				Optional:      true,
				Elem:          &pluginsdk.Schema{Type: pluginsdk.TypeString},
				Set:           pluginsdk.HashString,
				ConflictsWith: []string{"destination_port_range"},
			},

			"source_address_prefix": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ConflictsWith: []string{"source_address_prefixes"},
			},

			"source_address_prefixes": {
				Type:          pluginsdk.TypeSet,
				Optional:      true,
				Elem:          &pluginsdk.Schema{Type: pluginsdk.TypeString},
				Set:           pluginsdk.HashString,
				ConflictsWith: []string{"source_address_prefix"},
			},

			"destination_address_prefix": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ConflictsWith: []string{"destination_address_prefixes"},
			},

			"destination_address_prefixes": {
				Type:          pluginsdk.TypeSet,
				Optional:      true,
				Elem:          &pluginsdk.Schema{Type: pluginsdk.TypeString},
				Set:           pluginsdk.HashString,
				ConflictsWith: []string{"destination_address_prefix"},
			},

			// lintignore:S018
			"source_application_security_group_ids": {
				Type:     pluginsdk.TypeSet,
				MaxItems: 10,
				Optional: true,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
				Set:      pluginsdk.HashString,
			},

			// lintignore:S018
			"destination_application_security_group_ids": {
				Type:     pluginsdk.TypeSet,
				MaxItems: 10,
				Optional: true,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
				Set:      pluginsdk.HashString,
			},

			"access": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.SecurityRuleAccessAllow),
					string(network.SecurityRuleAccessDeny),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"priority": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(100, 4096),
			},

			"direction": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.SecurityRuleDirectionInbound),
					string(network.SecurityRuleDirectionOutbound),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},
		},
	}
}

func resourceNetworkSecurityRuleCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.SecurityRuleClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	nsgName := d.Get("network_security_group_name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, nsgName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Rule %q (Network Security Group %q / Resource Group %q): %s", name, nsgName, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_network_security_rule", *existing.ID)
		}
	}

	sourcePortRange := d.Get("source_port_range").(string)
	destinationPortRange := d.Get("destination_port_range").(string)
	sourceAddressPrefix := d.Get("source_address_prefix").(string)
	destinationAddressPrefix := d.Get("destination_address_prefix").(string)
	priority := int32(d.Get("priority").(int))
	access := d.Get("access").(string)
	direction := d.Get("direction").(string)
	protocol := d.Get("protocol").(string)

	if !meta.(*clients.Client).Features.Network.RelaxedLocking {
		locks.ByName(nsgName, networkSecurityGroupResourceName)
		defer locks.UnlockByName(nsgName, networkSecurityGroupResourceName)
	}

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
		r := r.(*pluginsdk.Set).List()
		for _, v := range r {
			s := v.(string)
			sourcePortRanges = append(sourcePortRanges, s)
		}
		rule.SecurityRulePropertiesFormat.SourcePortRanges = &sourcePortRanges
	}

	if r, ok := d.GetOk("destination_port_ranges"); ok {
		var destinationPortRanges []string
		r := r.(*pluginsdk.Set).List()
		for _, v := range r {
			s := v.(string)
			destinationPortRanges = append(destinationPortRanges, s)
		}
		rule.SecurityRulePropertiesFormat.DestinationPortRanges = &destinationPortRanges
	}

	if r, ok := d.GetOk("source_address_prefixes"); ok {
		var sourceAddressPrefixes []string
		r := r.(*pluginsdk.Set).List()
		for _, v := range r {
			s := v.(string)
			sourceAddressPrefixes = append(sourceAddressPrefixes, s)
		}
		rule.SecurityRulePropertiesFormat.SourceAddressPrefixes = &sourceAddressPrefixes
	}

	if r, ok := d.GetOk("destination_address_prefixes"); ok {
		var destinationAddressPrefixes []string
		r := r.(*pluginsdk.Set).List()
		for _, v := range r {
			s := v.(string)
			destinationAddressPrefixes = append(destinationAddressPrefixes, s)
		}
		rule.SecurityRulePropertiesFormat.DestinationAddressPrefixes = &destinationAddressPrefixes
	}

	if r, ok := d.GetOk("source_application_security_group_ids"); ok {
		var sourceApplicationSecurityGroups []network.ApplicationSecurityGroup
		for _, v := range r.(*pluginsdk.Set).List() {
			sg := network.ApplicationSecurityGroup{
				ID: utils.String(v.(string)),
			}
			sourceApplicationSecurityGroups = append(sourceApplicationSecurityGroups, sg)
		}
		rule.SourceApplicationSecurityGroups = &sourceApplicationSecurityGroups
	}

	if r, ok := d.GetOk("destination_application_security_group_ids"); ok {
		var destinationApplicationSecurityGroups []network.ApplicationSecurityGroup
		for _, v := range r.(*pluginsdk.Set).List() {
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

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
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

	return resourceNetworkSecurityRuleRead(d, meta)
}

func resourceNetworkSecurityRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.SecurityRuleClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
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

		if err := d.Set("source_application_security_group_ids", flattenApplicationSecurityGroupIds(props.SourceApplicationSecurityGroups)); err != nil {
			return fmt.Errorf("Error setting `source_application_security_group_ids`: %+v", err)
		}

		if err := d.Set("destination_application_security_group_ids", flattenApplicationSecurityGroupIds(props.DestinationApplicationSecurityGroups)); err != nil {
			return fmt.Errorf("Error setting `source_application_security_group_ids`: %+v", err)
		}
	}

	return nil
}

func resourceNetworkSecurityRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.SecurityRuleClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	nsgName := id.Path["networkSecurityGroups"]
	sgRuleName := id.Path["securityRules"]

	if !meta.(*clients.Client).Features.Network.RelaxedLocking {
		locks.ByName(nsgName, networkSecurityGroupResourceName)
		defer locks.UnlockByName(nsgName, networkSecurityGroupResourceName)
	}

	future, err := client.Delete(ctx, resGroup, nsgName, sgRuleName)
	if err != nil {
		return fmt.Errorf("Error Deleting Network Security Rule %q (NSG %q / Resource Group %q): %+v", sgRuleName, nsgName, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the deletion of Network Security Rule %q (NSG %q / Resource Group %q): %+v", sgRuleName, nsgName, resGroup, err)
	}

	return nil
}

func flattenApplicationSecurityGroupIds(groups *[]network.ApplicationSecurityGroup) []string {
	ids := make([]string, 0)

	if groups != nil {
		for _, v := range *groups {
			ids = append(ids, *v.ID)
		}
	}

	return ids
}
