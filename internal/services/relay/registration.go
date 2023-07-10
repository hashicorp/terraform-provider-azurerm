// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package relay

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var _ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}

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

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_relay_hybrid_connection":                    resourceArmRelayHybridConnection(),
		"azurerm_relay_hybrid_connection_authorization_rule": resourceRelayHybridConnectionAuthorizationRule(),
		"azurerm_relay_namespace":                            resourceRelayNamespace(),
		"azurerm_relay_namespace_authorization_rule":         resourceRelayNamespaceAuthorizationRule(),
	}
}
