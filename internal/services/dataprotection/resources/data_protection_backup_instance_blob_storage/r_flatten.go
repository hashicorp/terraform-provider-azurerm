// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_blob_storage

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/backupinstanceresources"
)

func flattenBlobBackupAutoProtection(input backupinstanceresources.BlobBackupDatasourceParametersForAutoProtection) (enabled bool, excludedContainerNamePrefixes []string) {
	excludedContainerNamePrefixes = make([]string, 0)
	settings := input.AutoProtectionSettings

	for _, rule := range pointer.From(settings.Rules) {
		// `Exclude` + `Prefix` is the only combination the API supports today, anything else is unknown to this resource
		if rule.Mode == backupinstanceresources.BlobBackupRuleModeExclude && rule.Type == backupinstanceresources.BlobBackupPatternTypePrefix {
			excludedContainerNamePrefixes = append(excludedContainerNamePrefixes, rule.Pattern)
		}
	}

	return settings.Enabled, excludedContainerNamePrefixes
}
