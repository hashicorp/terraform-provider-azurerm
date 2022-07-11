package dataprotection

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var _ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/data-protection"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "DataProtection"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"DataProtection",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_data_protection_backup_vault": dataSourceDataProtectionBackupVault(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_data_protection_backup_instance_blob_storage": resourceDataProtectionBackupInstanceBlobStorage(),
		"azurerm_data_protection_backup_instance_disk":         resourceDataProtectionBackupInstanceDisk(),
		"azurerm_data_protection_backup_instance_postgresql":   resourceDataProtectionBackupInstancePostgreSQL(),
		"azurerm_data_protection_backup_policy_blob_storage":   resourceDataProtectionBackupPolicyBlobStorage(),
		"azurerm_data_protection_backup_policy_disk":           resourceDataProtectionBackupPolicyDisk(),
		"azurerm_data_protection_backup_policy_postgresql":     resourceDataProtectionBackupPolicyPostgreSQL(),
		"azurerm_data_protection_backup_vault":                 resourceDataProtectionBackupVault(),
		"azurerm_data_protection_resource_guard":               resourceDataProtectionResourceGuard(),
	}
}
