package storageaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BlobRestoreStatus struct {
	FailureReason *string                    `json:"failureReason,omitempty"`
	Parameters    *BlobRestoreParameters     `json:"parameters,omitempty"`
	RestoreId     *string                    `json:"restoreId,omitempty"`
	Status        *BlobRestoreProgressStatus `json:"status,omitempty"`
}
