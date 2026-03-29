// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package relay

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration             = Registration{}
	_ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/relay"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Relay"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Messaging",
	}
}

// DataSources returns the supported Data Sources supported by this Service
func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

// Resources returns the supported Resources supported by this Service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		RelayNamespaceResource{},
		ArmRelayHybridConnectionResource{},
	}

	// return map[string]*pluginsdk.Resource{
	// 	"azurerm_relay_hybrid_connection":                    resourceArmRelayHybridConnection(),
	// 	"azurerm_relay_hybrid_connection_authorization_rule": resourceRelayHybridConnectionAuthorizationRule(),
	// 	"azurerm_relay_namespace":                            resourceRelayNamespace(),
	// 	"azurerm_relay_namespace_authorization_rule":         resourceRelayNamespaceAuthorizationRule(),
	// }
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
