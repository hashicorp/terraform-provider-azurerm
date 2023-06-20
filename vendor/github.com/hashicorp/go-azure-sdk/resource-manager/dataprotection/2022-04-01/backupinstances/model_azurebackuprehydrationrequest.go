package backupinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureBackupRehydrationRequest struct {
	RecoveryPointId              string               `json:"recoveryPointId"`
	RehydrationPriority          *RehydrationPriority `json:"rehydrationPriority,omitempty"`
	RehydrationRetentionDuration string               `json:"rehydrationRetentionDuration"`
}
