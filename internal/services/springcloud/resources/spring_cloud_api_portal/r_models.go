// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_api_portal

type SpringCloudAPIPortalModel struct {
	Name                       string              `tfschema:"name"`
	SpringCloudServiceId       string              `tfschema:"spring_cloud_service_id"`
	GatewayIds                 []string            `tfschema:"gateway_ids"`
	HttpsOnlyEnabled           bool                `tfschema:"https_only_enabled"`
	InstanceCount              int64               `tfschema:"instance_count"`
	PublicNetworkAccessEnabled bool                `tfschema:"public_network_access_enabled"`
	ApiTryOutEnabled           bool                `tfschema:"api_try_out_enabled"`
	Sso                        []ApiPortalSsoModel `tfschema:"sso"`
	Url                        string              `tfschema:"url"`
}

type ApiPortalSsoModel struct {
	ClientId     string   `tfschema:"client_id"`
	ClientSecret string   `tfschema:"client_secret"`
	IssuerUri    string   `tfschema:"issuer_uri"`
	Scope        []string `tfschema:"scope"`
}

func (s Resource) ModelObject() any {
	return &SpringCloudAPIPortalModel{}
}
