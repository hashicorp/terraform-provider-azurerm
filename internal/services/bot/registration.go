// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package bot

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}
	_ sdk.UntypedServiceRegistration               = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/bots"
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		AzureBotServiceResource{},
	}
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Bot"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Bot",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_bot_channel_alexa":              resourceBotChannelAlexa(),
		"azurerm_bot_channel_directline":         resourceBotChannelDirectline(),
		"azurerm_bot_channel_direct_line_speech": resourceBotChannelDirectLineSpeech(),
		"azurerm_bot_channel_email":              resourceBotChannelEmail(),
		"azurerm_bot_channel_facebook":           resourceBotChannelFacebook(),
		"azurerm_bot_channel_line":               resourceBotChannelLine(),
		"azurerm_bot_channel_ms_teams":           resourceBotChannelMsTeams(),
		"azurerm_bot_channel_slack":              resourceBotChannelSlack(),
		"azurerm_bot_channel_sms":                resourceBotChannelSMS(),
		"azurerm_bot_channel_web_chat":           resourceBotChannelWebChat(),
		"azurerm_bot_channels_registration":      resourceBotChannelsRegistration(),
		"azurerm_bot_connection":                 resourceArmBotConnection(),
		"azurerm_healthbot":                      resourceHealthbotService(),
		"azurerm_bot_web_app":                    resourceBotWebApp(),
	}
}
