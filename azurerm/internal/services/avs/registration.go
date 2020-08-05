package avs

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
    return "Avs"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
    return []string{
        "Avs",
    }
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
    return map[string]*schema.Resource{
        "azurerm_avs_private_cloud":    dataSourceAvsPrivateCloud(),
        "azurerm_avs_cluster":    dataSourceAvsCluster(),
        "azurerm_avs_hcx_enterprise_site":    dataSourceAvsHcxEnterpriseSite(),
        "azurerm_avs_authorization":    dataSourceAvsAuthorization(),
    }
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
    return map[string]*schema.Resource{
        "azurerm_avs_private_cloud":    resourceArmAvsPrivateCloud(),
        "azurerm_avs_cluster":    resourceArmAvsCluster(),
        "azurerm_avs_hcx_enterprise_site":    resourceArmAvsHcxEnterpriseSite(),
        "azurerm_avs_authorization":    resourceArmAvsAuthorization(),
    }
}
