// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_environment

type ApiCenterEnvironmentResourceModel struct {
	Name          string `tfschema:"name"`
	ServiceId     string `tfschema:"service_id"`
	Title         string `tfschema:"display_name"`
	Description   string `tfschema:"description"`
	Type          string `tfschema:"environment_type"`
	DevPortalUri  string `tfschema:"development_portal_uri"`
	Instructions  string `tfschema:"instructions"`
	ServerType    string `tfschema:"server_type"`
	MgmtPortalUri string `tfschema:"management_portal_uri"`
}

func (r Resource) ModelObject() any {
	return &ApiCenterEnvironmentResourceModel{}
}
