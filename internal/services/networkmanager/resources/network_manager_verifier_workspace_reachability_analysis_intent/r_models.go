// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_verifier_workspace_reachability_analysis_intent

func (Resource) ModelObject() any {
	return &ManagerVerifierWorkspaceReachabilityAnalysisIntentResourceModel{}
}

type ManagerVerifierWorkspaceReachabilityAnalysisIntentResourceModel struct {
	Description           string                                                        `tfschema:"description"`
	DestinationResourceId string                                                        `tfschema:"destination_resource_id"`
	IpTraffic             []ManagerVerifierWorkspaceReachabilityAnalysisIntentIpTraffic `tfschema:"ip_traffic"`
	Name                  string                                                        `tfschema:"name"`
	SourceResourceId      string                                                        `tfschema:"source_resource_id"`
	VerifierWorkspaceId   string                                                        `tfschema:"verifier_workspace_id"`
}
