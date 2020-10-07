package recoveryservices

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_recovery_services_vault": dataSourceArmRecoveryServicesVault(),
		"azurerm_backup_policy_vm":        dataSourceArmBackupPolicyVm(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_backup_container_storage_account":           resourceArmBackupProtectionContainerStorageAccount(),
		"azurerm_backup_policy_file_share":                   resourceArmBackupProtectionPolicyFileShare(),
		"azurerm_backup_protected_file_share":                resourceArmBackupProtectedFileShare(),
		"azurerm_backup_protected_vm":                        resourceArmRecoveryServicesBackupProtectedVM(),
		"azurerm_backup_policy_vm":                           resourceArmBackupProtectionPolicyVM(),
		"azurerm_recovery_services_vault":                    resourceArmRecoveryServicesVault(),
		"azurerm_site_recovery_fabric":                       resourceArmSiteRecoveryFabric(),
		"azurerm_site_recovery_network_mapping":              resourceArmSiteRecoveryNetworkMapping(),
		"azurerm_site_recovery_protection_container":         resourceArmSiteRecoveryProtectionContainer(),
		"azurerm_site_recovery_protection_container_mapping": resourceArmSiteRecoveryProtectionContainerMapping(),
		"azurerm_site_recovery_replicated_vm":                resourceArmSiteRecoveryReplicatedVM(),
		"azurerm_site_recovery_replication_policy":           resourceArmSiteRecoveryReplicationPolicy(),
	}
}
