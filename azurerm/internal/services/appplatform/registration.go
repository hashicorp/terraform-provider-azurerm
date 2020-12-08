package appplatform

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "App Platform"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Spring Cloud",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_spring_cloud_service": dataSourceArmSpringCloudService(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_spring_cloud_app":         resourceArmSpringCloudApp(),
		"azurerm_spring_cloud_certificate": resourceArmSpringCloudCertificate(),
		"azurerm_spring_cloud_service":     resourceArmSpringCloudService(),
	}
}
