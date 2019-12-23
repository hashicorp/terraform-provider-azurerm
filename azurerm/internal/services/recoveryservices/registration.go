package recoveryservices

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Recovery Services"
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_recovery_services_vault":                dataSourceArmRecoveryServicesVault(),
		"azurerm_recovery_services_protection_policy_vm": dataSourceArmRecoveryServicesProtectionPolicyVm(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_backup_container_storage_account":               resourceArmBackupProtectionContainerStorageAccount(),
		"azurerm_backup_policy_file_share":                       resourceArmBackupProtectionPolicyFileShare(),
		"azurerm_backup_protected_file_share":                    resourceArmBackupProtectedFileShare(),
		"azurerm_backup_protected_vm":                            resourceArmRecoveryServicesBackupProtectedVM(),
		"azurerm_backup_policy_vm":                               resourceArmBackupProtectionPolicyVM(),
		"azurerm_recovery_network_mapping":                       resourceArmRecoveryServicesNetworkMapping(),
		"azurerm_recovery_replicated_vm":                         resourceArmRecoveryServicesReplicatedVm(),
		"azurerm_recovery_services_fabric":                       resourceArmRecoveryServicesFabric(),
		"azurerm_recovery_services_protected_vm":                 resourceArmRecoveryServicesProtectedVm(),
		"azurerm_recovery_services_protection_container":         resourceArmRecoveryServicesProtectionContainer(),
		"azurerm_recovery_services_protection_container_mapping": resourceArmRecoveryServicesProtectionContainerMapping(),
		"azurerm_recovery_services_protection_policy_vm":         resourceArmRecoveryServicesProtectionPolicyVm(),
		"azurerm_recovery_services_replication_policy":           resourceArmRecoveryServicesReplicationPolicy(),
		"azurerm_recovery_services_vault":                        resourceArmRecoveryServicesVault(),
		"azurerm_site_recovery_fabric":                           resourceArmSiteRecoveryFabric(),
		"azurerm_site_recovery_network_mapping":                  resourceArmSiteRecoveryNetworkMapping(),
		"azurerm_site_recovery_protection_container":             resourceArmSiteRecoveryProtectionContainer(),
		"azurerm_site_recovery_protection_container_mapping":     resourceArmSiteRecoveryProtectionContainerMapping(),
		"azurerm_site_recovery_replicated_vm":                    resourceArmSiteRecoveryReplicatedVM(),
		"azurerm_site_recovery_replication_policy":               resourceArmSiteRecoveryReplicationPolicy(),
	}
}
