package signalr

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var _ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/signalr"
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
		"azurerm_signalr_service": dataSourceArmSignalRService(),
		"azurerm_web_pubsub":      dataSourceWebPubsub(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_signalr_service":             resourceArmSignalRService(),
		"azurerm_signalr_service_network_acl": resourceArmSignalRServiceNetworkACL(),
		"azurerm_web_pubsub":                  resourceWebPubSub(),
		"azurerm_web_pubsub_hub":              resourceWebPubsubHub(),
		"azurerm_web_pubsub_network_acl":      resourceWebpubsubNetworkACL(),
	}
}
