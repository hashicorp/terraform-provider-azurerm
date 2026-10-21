// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_admin_rule_collection

type ManagerAdminRuleCollectionModel struct {
	Name                         string   `tfschema:"name"`
	SecurityAdminConfigurationId string   `tfschema:"security_admin_configuration_id"`
	NetworkGroupIds              []string `tfschema:"network_group_ids"`
	Description                  string   `tfschema:"description"`
}

func (r Resource) ModelObject() any {
	return &ManagerAdminRuleCollectionModel{}
}
