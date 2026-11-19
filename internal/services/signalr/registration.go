// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package signalr

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/resources/signalr_service"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/resources/signalr_service_custom_certificate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/resources/signalr_service_custom_domain"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/resources/signalr_service_network_acl"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/resources/signalr_shared_private_link_resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/resources/web_pubsub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/resources/web_pubsub_custom_certificate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/resources/web_pubsub_custom_domain"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/resources/web_pubsub_hub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/resources/web_pubsub_network_acl"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/resources/web_pubsub_private_link_resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/resources/web_pubsub_shared_private_link_resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/resources/web_pubsub_socketio"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
	_ sdk.TypedServiceRegistrationWithAGitHubLabel   = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/signalr"
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		signalr_service_custom_certificate.Resource{},
		web_pubsub_custom_certificate.Resource{},
		signalr_service_custom_domain.Resource{},
		web_pubsub_custom_domain.Resource{},
		web_pubsub_socketio.Resource{},
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "SignalR"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Messaging",
		"Web PubSub",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_signalr_service":                  signalr_service.RegisterDataSource(),
		"azurerm_web_pubsub":                       web_pubsub.RegisterDataSource(),
		"azurerm_web_pubsub_private_link_resource": web_pubsub_private_link_resource.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_signalr_service":                         signalr_service.RegisterResource(),
		"azurerm_signalr_service_network_acl":             signalr_service_network_acl.RegisterResource(),
		"azurerm_signalr_shared_private_link_resource":    signalr_shared_private_link_resource.RegisterResource(),
		"azurerm_web_pubsub":                              web_pubsub.RegisterResource(),
		"azurerm_web_pubsub_hub":                          web_pubsub_hub.RegisterResource(),
		"azurerm_web_pubsub_network_acl":                  web_pubsub_network_acl.RegisterResource(),
		"azurerm_web_pubsub_shared_private_link_resource": web_pubsub_shared_private_link_resource.RegisterResource(),
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
	return []sdk.FrameworkListWrappedResource{
		signalr_service_custom_certificate.ListResource{},
		web_pubsub_custom_certificate.ListResource{},
		signalr_service_custom_domain.ListResource{},
		web_pubsub_custom_domain.ListResource{},
		signalr_service.ListResource{},
		web_pubsub.ListResource{},
		web_pubsub_socketio.ListResource{},
	}
}
