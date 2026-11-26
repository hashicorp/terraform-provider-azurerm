// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package relay_namespace_authorization_rule

type RelayNamespaceAuthorizationRuleResourceModel struct {
	Name                      string `tfschema:"name"`
	NamespaceName             string `tfschema:"namespace_name"`
	ResourceGroupName         string `tfschema:"resource_group_name"`
	Listen                    bool   `tfschema:"listen"`
	Manage                    bool   `tfschema:"manage"`
	Send                      bool   `tfschema:"send"`
	PrimaryConnectionString   string `tfschema:"primary_connection_string"`
	PrimaryKey                string `tfschema:"primary_key"`
	SecondaryConnectionString string `tfschema:"secondary_connection_string"`
	SecondaryKey              string `tfschema:"secondary_key"`
}

func (r Resource) ModelObject() any {
	return &RelayNamespaceAuthorizationRuleResourceModel{}
}
