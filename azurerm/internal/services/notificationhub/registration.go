package notificationhub

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Notification Hub"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Messaging",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_notification_hub_namespace": dataSourceNotificationHubNamespace(),
		"azurerm_notification_hub":           dataSourceNotificationHub(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_notification_hub_authorization_rule": resourceNotificationHubAuthorizationRule(),
		"azurerm_notification_hub_namespace":          resourceNotificationHubNamespace(),
		"azurerm_notification_hub":                    resourceNotificationHub(),
	}
}
