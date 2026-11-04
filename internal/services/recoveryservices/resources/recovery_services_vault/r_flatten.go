// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package recovery_services_vault

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2025-08-01/vaults"
)

func flattenVaultEncryption(model vaults.Vault) any {
	if model.Properties == nil || model.Properties.Encryption == nil {
		return nil
	}
	encryption := model.Properties.Encryption
	if encryption.KeyVaultProperties == nil || encryption.KeyVaultProperties.KeyUri == nil {
		return nil
	}
	if encryption.KekIdentity == nil || encryption.KekIdentity.UseSystemAssignedIdentity == nil {
		return nil
	}
	encryptionMap := make(map[string]any)
	encryptionMap["key_id"] = encryption.KeyVaultProperties.KeyUri
	encryptionMap["use_system_assigned_identity"] = *encryption.KekIdentity.UseSystemAssignedIdentity
	encryptionMap["infrastructure_encryption_enabled"] = *encryption.InfrastructureEncryption == vaults.InfrastructureEncryptionStateEnabled
	if encryption.KekIdentity.UserAssignedIdentity != nil {
		encryptionMap["user_assigned_identity_id"] = *encryption.KekIdentity.UserAssignedIdentity
	}
	return encryptionMap
}

func flattenRecoveryServicesVaultPublicNetworkAccess(input *vaults.PublicNetworkAccess) bool {
	if input == nil {
		return false
	}
	return *input == vaults.PublicNetworkAccessEnabled
}

func flattenRecoveryServicesVaultMonitorSettings(input *vaults.MonitoringSettings) []any {
	// `Monitor` is an optional parameters, and won't be returned from API if it has not been specified.
	if input == nil {
		return []any{}
	}
	allJobAlert := false
	allFailoverAlert := false
	allReplicationAlert := false
	criticalAlert := false
	emailNotification := false

	if input.AzureMonitorAlertSettings != nil {
		allJobAlert = pointer.From(input.AzureMonitorAlertSettings.AlertsForAllJobFailures) == vaults.AlertsStateEnabled
		allFailoverAlert = pointer.From(input.AzureMonitorAlertSettings.AlertsForAllFailoverIssues) == vaults.AlertsStateEnabled
		allReplicationAlert = pointer.From(input.AzureMonitorAlertSettings.AlertsForAllReplicationIssues) == vaults.AlertsStateEnabled
	}

	if input.ClassicAlertSettings != nil {
		criticalAlert = pointer.From(input.ClassicAlertSettings.AlertsForCriticalOperations) == vaults.AlertsStateEnabled
		emailNotification = pointer.From(input.ClassicAlertSettings.EmailNotificationsForSiteRecovery) == vaults.AlertsStateEnabled
	}

	return []any{
		map[string]any{
			"alerts_for_all_job_failures_enabled":            allJobAlert,
			"alerts_for_all_failover_issues_enabled":         allFailoverAlert,
			"alerts_for_all_replication_issues_enabled":      allReplicationAlert,
			"alerts_for_critical_operation_failures_enabled": criticalAlert,
			"email_notifications_for_site_recovery_enabled":  emailNotification,
		},
	}
}
