// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bot_service_azure_bot

type BotServiceModel struct {
	Name                              string         `tfschema:"name"`
	ResourceGroupName                 string         `tfschema:"resource_group_name"`
	Location                          string         `tfschema:"location"`
	Sku                               string         `tfschema:"sku"`
	MicrosoftAppId                    string         `tfschema:"microsoft_app_id"`
	DisplayName                       string         `tfschema:"display_name"`
	Endpoint                          string         `tfschema:"endpoint"`
	DeveloperAppInsightsKey           string         `tfschema:"developer_app_insights_key"`
	DeveloperAppInsightsApiKey        string         `tfschema:"developer_app_insights_api_key"`
	DeveloperAppInsightsApplicationId string         `tfschema:"developer_app_insights_application_id"`
	CmkKeyVaultKeyUrl                 string         `tfschema:"cmk_key_vault_key_url"`
	MicrosoftAppMsiId                 string         `tfschema:"microsoft_app_msi_id"`
	MicrosoftAppTenantId              string         `tfschema:"microsoft_app_tenant_id"`
	MicrosoftAppType                  string         `tfschema:"microsoft_app_type"`
	LocalAuthenticationEnabled        bool           `tfschema:"local_authentication_enabled"`
	LuisAppIds                        []string       `tfschema:"luis_app_ids"`
	LuisKey                           string         `tfschema:"luis_key"`
	PublicNetworkAccessEnabled        bool           `tfschema:"public_network_access_enabled"`
	StreamingEndpointEnabled          bool           `tfschema:"streaming_endpoint_enabled"`
	IconUrl                           string         `tfschema:"icon_url"`
	Tags                              map[string]any `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &BotServiceModel{}
}
