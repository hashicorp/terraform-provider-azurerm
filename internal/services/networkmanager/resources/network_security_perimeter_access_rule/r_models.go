// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_security_perimeter_access_rule

type NetworkSecurityPerimeterAccessRuleResourceModel struct {
	Name                      string   `tfschema:"name"`
	ProfileId                 string   `tfschema:"network_security_perimeter_profile_id"`
	Direction                 string   `tfschema:"direction"`
	AddressPrefixes           []string `tfschema:"address_prefixes"`
	FullyQualifiedDomainNames []string `tfschema:"fqdns"`
	Subscriptions             []string `tfschema:"subscription_ids"`
	ServiceTags               []string `tfschema:"service_tags"`
}

func (Resource) ModelObject() any {
	return &NetworkSecurityPerimeterAccessRuleResourceModel{}
}
