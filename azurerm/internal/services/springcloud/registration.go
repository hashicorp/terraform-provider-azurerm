package springcloud

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Spring Cloud"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Spring Cloud",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_spring_cloud_app":     dataSourceSpringCloudApp(),
		"azurerm_spring_cloud_service": dataSourceSpringCloudService(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_spring_cloud_active_deployment":        resourceSpringCloudActiveDeployment(),
		"azurerm_spring_cloud_app":                      resourceSpringCloudApp(),
		"azurerm_spring_cloud_app_cosmosdb_association": resourceSpringCloudAppCosmosDBAssociation(),
		"azurerm_spring_cloud_app_mysql_association":    resourceSpringCloudAppMysqlAssociation(),
		"azurerm_spring_cloud_app_redis_association":    resourceSpringCloudAppRedisAssociation(),
		"azurerm_spring_cloud_certificate":              resourceSpringCloudCertificate(),
		"azurerm_spring_cloud_custom_domain":            resourceSpringCloudCustomDomain(),
		"azurerm_spring_cloud_java_deployment":          resourceSpringCloudJavaDeployment(),
		"azurerm_spring_cloud_service":                  resourceSpringCloudService(),
	}
}
