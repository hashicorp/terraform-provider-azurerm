// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_edge_action_version

type EdgeActionVersionModel struct {
	Name                     string         `tfschema:"name"`
	CdnFrontDoorEdgeActionId string         `tfschema:"cdn_frontdoor_edge_action_id"`
	DeploymentType           string         `tfschema:"deployment_type"`
	CodeContentBase64        string         `tfschema:"code_content_base64"`
	CodeFileName             string         `tfschema:"code_file_name"`
	Tags                     map[string]any `tfschema:"tags"`
	DefaultVersionEnabled    bool           `tfschema:"default_version_enabled"`
	ValidationStatus         string         `tfschema:"validation_status"`
}

func (Resource) ModelObject() any {
	return &EdgeActionVersionModel{}
}
