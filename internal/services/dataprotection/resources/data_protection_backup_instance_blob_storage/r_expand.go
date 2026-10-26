// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_blob_storage

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/backupinstanceresources"
)

func expandBlobBackupAutoProtection(excludedContainerNamePrefixes []any) backupinstanceresources.BlobBackupDatasourceParametersForAutoProtection {
	settings := backupinstanceresources.BlobBackupRuleBasedAutoProtectionSettings{
		Enabled: true,
	}

	// without any rules every present and future container is eligible for auto protection
	if len(excludedContainerNamePrefixes) > 0 {
		rules := make([]backupinstanceresources.BlobBackupAutoProtectionRule, 0, len(excludedContainerNamePrefixes))
		for _, prefix := range excludedContainerNamePrefixes {
			rules = append(rules, backupinstanceresources.BlobBackupAutoProtectionRule{
				// `objectType` is a required plain field on the rule rather than a discriminator, so the SDK does not set it
				ObjectType: "BlobBackupAutoProtectionRule",
				Mode:       backupinstanceresources.BlobBackupRuleModeExclude,
				Type:       backupinstanceresources.BlobBackupPatternTypePrefix,
				Pattern:    prefix.(string),
			})
		}
		settings.Rules = &rules
	}

	return backupinstanceresources.BlobBackupDatasourceParametersForAutoProtection{
		AutoProtectionSettings: settings,
	}
}
