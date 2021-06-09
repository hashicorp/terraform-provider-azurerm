package recoveryservices

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Recovery Services"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Recovery Services",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_recovery_services_vault": dataSourceRecoveryServicesVault(),
		"azurerm_backup_policy_vm":        dataSourceBackupPolicyVm(),
	}
}

// todo - this package should probably be split into backup, recovery, and site recovery?
// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_backup_container_storage_account":           resourceBackupProtectionContainerStorageAccount(),
		"azurerm_backup_policy_file_share":                   resourceBackupProtectionPolicyFileShare(),
		"azurerm_backup_protected_file_share":                resourceBackupProtectedFileShare(),
		"azurerm_backup_protected_vm":                        resourceRecoveryServicesBackupProtectedVM(),
		"azurerm_backup_policy_vm":                           resourceBackupProtectionPolicyVM(),
		"azurerm_recovery_services_vault":                    resourceRecoveryServicesVault(),
		"azurerm_site_recovery_fabric":                       resourceSiteRecoveryFabric(),
		"azurerm_site_recovery_network_mapping":              resourceSiteRecoveryNetworkMapping(),
		"azurerm_site_recovery_protection_container":         resourceSiteRecoveryProtectionContainer(),
		"azurerm_site_recovery_protection_container_mapping": resourceSiteRecoveryProtectionContainerMapping(),
		"azurerm_site_recovery_replicated_vm":                resourceSiteRecoveryReplicatedVM(),
		"azurerm_site_recovery_replication_policy":           resourceSiteRecoveryReplicationPolicy(),
	}
}
