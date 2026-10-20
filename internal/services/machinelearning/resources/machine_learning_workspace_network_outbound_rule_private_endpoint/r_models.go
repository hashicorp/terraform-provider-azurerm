// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machine_learning_workspace_network_outbound_rule_private_endpoint

type machineLearningWorkspaceOutboundRulePrivateEndpointModel struct {
	Name              string `tfschema:"name"`
	WorkspaceId       string `tfschema:"workspace_id"`
	ServiceResourceId string `tfschema:"service_resource_id"`
	SubresourceTarget string `tfschema:"sub_resource_target"`
	SparkEnabled      bool   `tfschema:"spark_enabled"`
}

func (r Resource) ModelObject() any {
	return &machineLearningWorkspaceOutboundRulePrivateEndpointModel{}
}
