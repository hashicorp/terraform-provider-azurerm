// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_admin_rule

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/adminrules"
)

type ManagerAdminRuleModel struct {
	Name                    string                                        `tfschema:"name"`
	NetworkRuleCollectionId string                                        `tfschema:"admin_rule_collection_id"`
	Action                  adminrules.SecurityConfigurationRuleAccess    `tfschema:"action"`
	Description             string                                        `tfschema:"description"`
	DestinationPortRanges   []string                                      `tfschema:"destination_port_ranges"`
	Destinations            []AddressPrefixItemModel                      `tfschema:"destination"`
	Direction               adminrules.SecurityConfigurationRuleDirection `tfschema:"direction"`
	Priority                int64                                         `tfschema:"priority"`
	Protocol                adminrules.SecurityConfigurationRuleProtocol  `tfschema:"protocol"`
	SourcePortRanges        []string                                      `tfschema:"source_port_ranges"`
	Sources                 []AddressPrefixItemModel                      `tfschema:"source"`
}

type AddressPrefixItemModel struct {
	AddressPrefix     string                       `tfschema:"address_prefix"`
	AddressPrefixType adminrules.AddressPrefixType `tfschema:"address_prefix_type"`
}

func (r Resource) ModelObject() any {
	return &ManagerAdminRuleModel{}
}
