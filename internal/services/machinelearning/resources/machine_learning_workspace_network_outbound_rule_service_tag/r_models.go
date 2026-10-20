// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machine_learning_workspace_network_outbound_rule_service_tag

type MachineLearningWorkspaceServiceTagOutboundRuleModel struct {
	Name        string `tfschema:"name"`
	WorkspaceId string `tfschema:"workspace_id"`
	ServiceTag  string `tfschema:"service_tag"`
	Protocol    string `tfschema:"protocol"`
	PortRanges  string `tfschema:"port_ranges"`
}

func (r Resource) ModelObject() any {
	return &MachineLearningWorkspaceServiceTagOutboundRuleModel{}
}
