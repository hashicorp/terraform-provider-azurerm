// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_security_group

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networksecuritygroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/set"
)

func resourceNetworkSecurityGroupFlatten(d *pluginsdk.ResourceData, id *networksecuritygroups.NetworkSecurityGroupId, nsg *networksecuritygroups.NetworkSecurityGroup) error {
	d.Set("name", id.NetworkSecurityGroupName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if nsg != nil {
		d.Set("location", location.NormalizeNilable(nsg.Location))
		if props := nsg.Properties; props != nil {
			if err := d.Set("security_rule", flattenNetworkSecurityRules(props.SecurityRules)); err != nil {
				return fmt.Errorf("setting `security_rule`: %+v", err)
			}
		}
		if err := tags.FlattenAndSet(d, nsg.Tags); err != nil {
			return err
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenNetworkSecurityRules(rules *[]networksecuritygroups.SecurityRule) []map[string]any {
	result := make([]map[string]any, 0)

	// For fixing the case insensitive issue for the NSR protocol in Azure
	// See: https://github.com/hashicorp/terraform-provider-azurerm/issues/16092
	protocolMap := map[string]string{}
	for _, protocol := range networksecuritygroups.PossibleValuesForSecurityRuleProtocol() {
		protocolMap[strings.ToLower(protocol)] = protocol
	}

	if rules != nil {
		for _, rule := range *rules {
			sgRule := make(map[string]any)
			sgRule["name"] = *rule.Name

			if props := rule.Properties; props != nil {
				if props.Description != nil {
					sgRule["description"] = *props.Description
				}

				if props.DestinationAddressPrefix != nil {
					sgRule["destination_address_prefix"] = *props.DestinationAddressPrefix
				}
				if props.DestinationAddressPrefixes != nil {
					sgRule["destination_address_prefixes"] = set.FromStringSlice(*props.DestinationAddressPrefixes)
				}
				if props.DestinationPortRange != nil {
					sgRule["destination_port_range"] = *props.DestinationPortRange
				}
				if props.DestinationPortRanges != nil {
					sgRule["destination_port_ranges"] = set.FromStringSlice(*props.DestinationPortRanges)
				}

				destinationApplicationSecurityGroups := make([]string, 0)
				if props.DestinationApplicationSecurityGroups != nil {
					for _, g := range *props.DestinationApplicationSecurityGroups {
						destinationApplicationSecurityGroups = append(destinationApplicationSecurityGroups, *g.Id)
					}
				}
				sgRule["destination_application_security_group_ids"] = set.FromStringSlice(destinationApplicationSecurityGroups)

				if props.SourceAddressPrefix != nil {
					sgRule["source_address_prefix"] = *props.SourceAddressPrefix
				}
				if props.SourceAddressPrefixes != nil {
					sgRule["source_address_prefixes"] = set.FromStringSlice(*props.SourceAddressPrefixes)
				}

				sourceApplicationSecurityGroups := make([]string, 0)
				if props.SourceApplicationSecurityGroups != nil {
					for _, g := range *props.SourceApplicationSecurityGroups {
						sourceApplicationSecurityGroups = append(sourceApplicationSecurityGroups, *g.Id)
					}
				}
				sgRule["source_application_security_group_ids"] = set.FromStringSlice(sourceApplicationSecurityGroups)

				if props.SourcePortRange != nil {
					sgRule["source_port_range"] = *props.SourcePortRange
				}
				if props.SourcePortRanges != nil {
					sgRule["source_port_ranges"] = set.FromStringSlice(*props.SourcePortRanges)
				}

				sgRule["protocol"] = protocolMap[strings.ToLower(string(props.Protocol))]
				sgRule["priority"] = int(props.Priority)
				sgRule["access"] = string(props.Access)
				sgRule["direction"] = string(props.Direction)
			}

			result = append(result, sgRule)
		}
	}

	return result
}
