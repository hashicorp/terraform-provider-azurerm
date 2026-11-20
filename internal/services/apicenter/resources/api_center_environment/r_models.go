// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_environment

type ApiCenterEnvironmentResourceModel struct {
	Name             string `tfschema:"name"`
	ServiceId        string `tfschema:"api_center_service_id"`
	Title            string `tfschema:"display_name"`
	Description      string `tfschema:"description"`
	Type             string `tfschema:"environment_type"`
	DevPortalUrl     string `tfschema:"development_portal_url"`
	Instructions     string `tfschema:"instructions"`
	ServerType       string `tfschema:"server_type"`
	MgmtPortalUrl    string `tfschema:"management_portal_url"`
	CustomProperties string `tfschema:"custom_properties_json"`
}

func (r Resource) ModelObject() any {
	return &ApiCenterEnvironmentResourceModel{}
}
