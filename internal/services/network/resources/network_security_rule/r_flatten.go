// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_security_rule

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/securityrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceNetworkSecurityRuleFlatten(d *pluginsdk.ResourceData, id *securityrules.SecurityRuleId, model *securityrules.SecurityRule) error {
	d.Set("name", id.SecurityRuleName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("network_security_group_name", id.NetworkSecurityGroupName)

	// For fixing the case insensitive issue for the NSR protocol in Azure
	// See: https://github.com/hashicorp/terraform-provider-azurerm/issues/16092
	protocolMap := map[string]securityrules.SecurityRuleProtocol{}
	for _, protocol := range securityrules.PossibleValuesForSecurityRuleProtocol() {
		protocolMap[strings.ToLower(protocol)] = securityrules.SecurityRuleProtocol(protocol)
	}

	if model != nil {
		if props := model.Properties; props != nil {
			d.Set("description", props.Description)
			d.Set("protocol", string(protocolMap[strings.ToLower(string(props.Protocol))]))
			d.Set("destination_address_prefix", props.DestinationAddressPrefix)
			d.Set("destination_address_prefixes", props.DestinationAddressPrefixes)
			d.Set("destination_port_range", props.DestinationPortRange)
			d.Set("destination_port_ranges", props.DestinationPortRanges)
			d.Set("source_address_prefix", props.SourceAddressPrefix)
			d.Set("source_address_prefixes", props.SourceAddressPrefixes)
			d.Set("source_port_range", props.SourcePortRange)
			d.Set("source_port_ranges", props.SourcePortRanges)
			d.Set("access", string(props.Access))
			d.Set("priority", int(props.Priority))
			d.Set("direction", string(props.Direction))

			if err := d.Set("source_application_security_group_ids", flattenApplicationSecurityGroupIds(props.SourceApplicationSecurityGroups)); err != nil {
				return fmt.Errorf("setting `source_application_security_group_ids`: %+v", err)
			}

			if err := d.Set("destination_application_security_group_ids", flattenApplicationSecurityGroupIds(props.DestinationApplicationSecurityGroups)); err != nil {
				return fmt.Errorf("setting `source_application_security_group_ids`: %+v", err)
			}
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenApplicationSecurityGroupIds(groups *[]securityrules.ApplicationSecurityGroup) []string {
	ids := make([]string, 0)

	if groups != nil {
		for _, v := range *groups {
			ids = append(ids, *v.Id)
		}
	}

	return ids
}
