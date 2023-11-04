package longtermretentionbackups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LongTermRetentionOperationResultProperties struct {
	FromBackupResourceId          *string                  `json:"fromBackupResourceId,omitempty"`
	Message                       *string                  `json:"message,omitempty"`
	OperationType                 *string                  `json:"operationType,omitempty"`
	RequestId                     *string                  `json:"requestId,omitempty"`
	Status                        *string                  `json:"status,omitempty"`
	TargetBackupStorageRedundancy *BackupStorageRedundancy `json:"targetBackupStorageRedundancy,omitempty"`
	ToBackupResourceId            *string                  `json:"toBackupResourceId,omitempty"`
}
