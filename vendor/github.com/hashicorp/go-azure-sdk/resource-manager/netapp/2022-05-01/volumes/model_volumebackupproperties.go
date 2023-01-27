package volumes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VolumeBackupProperties struct {
	BackupEnabled  *bool   `json:"backupEnabled,omitempty"`
	BackupPolicyId *string `json:"backupPolicyId,omitempty"`
	PolicyEnforced *bool   `json:"policyEnforced,omitempty"`
	VaultId        *string `json:"vaultId,omitempty"`
}
