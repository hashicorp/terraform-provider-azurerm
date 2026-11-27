// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_dev_tool_portal

type SpringCloudDevToolPortalModel struct {
	Name                          string     `tfschema:"name"`
	SpringCloudServiceId          string     `tfschema:"spring_cloud_service_id"`
	ApplicationAcceleratorEnabled bool       `tfschema:"application_accelerator_enabled"`
	ApplicationLiveViewEnabled    bool       `tfschema:"application_live_view_enabled"`
	PublicNetworkAccessEnabled    bool       `tfschema:"public_network_access_enabled"`
	Sso                           []SsoModel `tfschema:"sso"`
}

type SsoModel struct {
	ClientId     string   `tfschema:"client_id"`
	ClientSecret string   `tfschema:"client_secret"`
	MetadataUrl  string   `tfschema:"metadata_url"`
	Scope        []string `tfschema:"scope"`
}

func (s Resource) ModelObject() any {
	return &SpringCloudDevToolPortalModel{}
}
