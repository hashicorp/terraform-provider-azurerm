// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machine_learning_workspace_network_outbound_rule_fqdn

type machineLearningWorkspaceOutboundRuleFqdnModel struct {
	Name        string `tfschema:"name"`
	WorkspaceId string `tfschema:"workspace_id"`
	Destination string `tfschema:"destination_fqdn"`
}

func (r Resource) ModelObject() any {
	return &machineLearningWorkspaceOutboundRuleFqdnModel{}
}
