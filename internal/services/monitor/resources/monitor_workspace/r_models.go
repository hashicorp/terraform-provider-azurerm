// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_workspace

type WorkspaceResourceModel struct {
	Name                            string            `tfschema:"name"`
	ResourceGroupName               string            `tfschema:"resource_group_name"`
	QueryEndpoint                   string            `tfschema:"query_endpoint"`
	DefaultDataCollectionEndpointId string            `tfschema:"default_data_collection_endpoint_id"`
	DefaultDataCollectionRuleId     string            `tfschema:"default_data_collection_rule_id"`
	PublicNetworkAccessEnabled      bool              `tfschema:"public_network_access_enabled"`
	Location                        string            `tfschema:"location"`
	Tags                            map[string]string `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &WorkspaceResourceModel{}
}
