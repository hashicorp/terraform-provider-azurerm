// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bot

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/resources/bot_channel_alexa"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/resources/bot_channel_direct_line_speech"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/resources/bot_channel_directline"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/resources/bot_channel_email"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/resources/bot_channel_facebook"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/resources/bot_channel_line"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/resources/bot_channel_ms_teams"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/resources/bot_channel_slack"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/resources/bot_channel_sms"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/resources/bot_channel_web_chat"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/resources/bot_channels_registration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/resources/bot_connection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/resources/bot_service_azure_bot"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/resources/bot_web_app"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/resources/healthbot"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration             = Registration{}
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
		bot_service_azure_bot.Resource{},
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
		"azurerm_bot_channel_alexa":              bot_channel_alexa.RegisterResource(),
		"azurerm_bot_channel_direct_line_speech": bot_channel_direct_line_speech.RegisterResource(),
		"azurerm_bot_channel_directline":         bot_channel_directline.RegisterResource(),
		"azurerm_bot_channel_email":              bot_channel_email.RegisterResource(),
		"azurerm_bot_channel_facebook":           bot_channel_facebook.RegisterResource(),
		"azurerm_bot_channel_line":               bot_channel_line.RegisterResource(),
		"azurerm_bot_channel_ms_teams":           bot_channel_ms_teams.RegisterResource(),
		"azurerm_bot_channel_slack":              bot_channel_slack.RegisterResource(),
		"azurerm_bot_channel_sms":                bot_channel_sms.RegisterResource(),
		"azurerm_bot_channel_web_chat":           bot_channel_web_chat.RegisterResource(),
		"azurerm_bot_channels_registration":      bot_channels_registration.RegisterResource(),
		"azurerm_bot_connection":                 bot_connection.RegisterResource(),
		"azurerm_bot_web_app":                    bot_web_app.RegisterResource(),
		"azurerm_healthbot":                      healthbot.RegisterResource(),
	}
}

func (r Registration) Actions() []func() action.Action {
	return []func() action.Action{}
}

func (r Registration) FrameworkResources() []sdk.FrameworkWrappedResource {
	return []sdk.FrameworkWrappedResource{}
}

func (r Registration) FrameworkDataSources() []sdk.FrameworkWrappedDataSource {
	return []sdk.FrameworkWrappedDataSource{}
}

func (r Registration) EphemeralResources() []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (r Registration) ListResources() []sdk.FrameworkListWrappedResource {
	return []sdk.FrameworkListWrappedResource{}
}
