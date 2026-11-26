// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package recovery_services_vault

import (
	"errors"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2025-08-01/vaults"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandEncryption(d *pluginsdk.ResourceData) (*vaults.VaultPropertiesEncryption, error) {
	encryptionRaw := d.Get("encryption")
	if encryptionRaw == nil {
		return nil, nil
	}
	settings := encryptionRaw.([]any)
	if len(settings) == 0 {
		return nil, nil
	}
	encryptionMap := settings[0].(map[string]any)
	keyUri := encryptionMap["key_id"].(string)
	enabledInfraEncryption := encryptionMap["infrastructure_encryption_enabled"].(bool)
	infraEncryptionState := vaults.InfrastructureEncryptionStateEnabled
	if !enabledInfraEncryption {
		infraEncryptionState = vaults.InfrastructureEncryptionStateDisabled
	}
	encryption := &vaults.VaultPropertiesEncryption{
		KeyVaultProperties: &vaults.CmkKeyVaultProperties{
			KeyUri: pointer.To(keyUri),
		},
		KekIdentity: &vaults.CmkKekIdentity{
			UseSystemAssignedIdentity: pointer.To(encryptionMap["use_system_assigned_identity"].(bool)),
		},
		InfrastructureEncryption: &infraEncryptionState,
	}
	if v, ok := encryptionMap["user_assigned_identity_id"].(string); ok && v != "" {
		if *encryption.KekIdentity.UseSystemAssignedIdentity {
			return nil, errors.New("`use_system_assigned_identity` must be disabled when `user_assigned_identity_id` is set")
		}
		encryption.KekIdentity.UserAssignedIdentity = pointer.To(v)
	}
	return encryption, nil
}

func expandRecoveryServicesVaultSecuritySettings(input any) *vaults.SecuritySettings {
	if input == nil || len(input.(string)) == 0 {
		return nil
	}
	return &vaults.SecuritySettings{
		ImmutabilitySettings: &vaults.ImmutabilitySettings{
			State: pointer.ToEnum[vaults.ImmutabilityState](input.(string)),
		},
	}
}

func expandRecoveryServicesVaultPublicNetworkAccess(input bool) *vaults.PublicNetworkAccess {
	out := vaults.PublicNetworkAccessDisabled
	if input {
		out = vaults.PublicNetworkAccessEnabled
	}
	return &out
}

func expandRecoveryServicesVaultMonitorSettings(input []any) *vaults.MonitoringSettings {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]any)

	allJobAlert := vaults.AlertsStateDisabled
	if v["alerts_for_all_job_failures_enabled"].(bool) {
		allJobAlert = vaults.AlertsStateEnabled
	}

	allFailoverAlert := vaults.AlertsStateDisabled
	if v["alerts_for_all_failover_issues_enabled"].(bool) {
		allFailoverAlert = vaults.AlertsStateEnabled
	}

	allReplicationAlert := vaults.AlertsStateDisabled
	if v["alerts_for_all_replication_issues_enabled"].(bool) {
		allReplicationAlert = vaults.AlertsStateEnabled
	}

	criticalOperation := vaults.AlertsStateDisabled
	if v["alerts_for_critical_operation_failures_enabled"].(bool) {
		criticalOperation = vaults.AlertsStateEnabled
	}

	emailNotification := vaults.AlertsStateDisabled
	if v["email_notifications_for_site_recovery_enabled"].(bool) {
		emailNotification = vaults.AlertsStateEnabled
	}

	return pointer.To(vaults.MonitoringSettings{
		AzureMonitorAlertSettings: pointer.To(vaults.AzureMonitorAlertSettings{
			AlertsForAllJobFailures:       pointer.To(allJobAlert),
			AlertsForAllFailoverIssues:    pointer.To(allFailoverAlert),
			AlertsForAllReplicationIssues: pointer.To(allReplicationAlert),
		}),
		ClassicAlertSettings: pointer.To(vaults.ClassicAlertSettings{
			AlertsForCriticalOperations:       pointer.To(criticalOperation),
			EmailNotificationsForSiteRecovery: pointer.To(emailNotification),
		}),
	})
}
