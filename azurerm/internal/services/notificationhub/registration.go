package notificationhub

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_notification_hub_namespace": dataSourceNotificationHubNamespace(),
		"azurerm_notification_hub":           dataSourceNotificationHub(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_notification_hub_authorization_rule": resourceArmNotificationHubAuthorizationRule(),
		"azurerm_notification_hub_namespace":          resourceArmNotificationHubNamespace(),
		"azurerm_notification_hub":                    resourceArmNotificationHub(),
	}
}
