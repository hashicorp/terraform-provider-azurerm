package bot

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

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
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_bot_channel_directline":    resourceArmBotChannelDirectline(),
		"azurerm_bot_channel_email":         resourceArmBotChannelEmail(),
		"azurerm_bot_channel_ms_teams":      resourceArmBotChannelMsTeams(),
		"azurerm_bot_channel_slack":         resourceArmBotChannelSlack(),
		"azurerm_bot_channels_registration": resourceArmBotChannelsRegistration(),
		"azurerm_bot_connection":            resourceArmBotConnection(),
		"azurerm_bot_web_app":               resourceArmBotWebApp(),
	}
}
