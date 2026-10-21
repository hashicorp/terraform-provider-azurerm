// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_security_admin_configuration

type ManagerSecurityAdminConfigurationModel struct {
	Name                                    string   `tfschema:"name"`
	NetworkManagerId                        string   `tfschema:"network_manager_id"`
	ApplyOnNetworkIntentPolicyBasedServices []string `tfschema:"apply_on_network_intent_policy_based_services"`
	Description                             string   `tfschema:"description"`
}

func (r Resource) ModelObject() any {
	return &ManagerSecurityAdminConfigurationModel{}
}
