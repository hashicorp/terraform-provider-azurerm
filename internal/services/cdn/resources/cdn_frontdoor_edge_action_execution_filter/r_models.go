// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_edge_action_execution_filter

type EdgeActionExecutionFilterModel struct {
	Name                            string         `tfschema:"name"`
	CdnFrontDoorEdgeActionId        string         `tfschema:"cdn_frontdoor_edge_action_id"`
	CdnFrontDoorEdgeActionVersionId string         `tfschema:"cdn_frontdoor_edge_action_version_id"`
	HeaderName                      string         `tfschema:"header_name"`
	HeaderValue                     string         `tfschema:"header_value"`
	Tags                            map[string]any `tfschema:"tags"`
}

func (Resource) ModelObject() any {
	return &EdgeActionExecutionFilterModel{}
}
