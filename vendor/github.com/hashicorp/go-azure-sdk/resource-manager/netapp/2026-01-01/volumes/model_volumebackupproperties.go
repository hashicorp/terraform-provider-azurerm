package volumes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VolumeBackupProperties struct {
	BackupPolicyId *string `json:"backupPolicyId,omitempty"`
	BackupVaultId  *string `json:"backupVaultId,omitempty"`
	PolicyEnforced *bool   `json:"policyEnforced,omitempty"`
}
