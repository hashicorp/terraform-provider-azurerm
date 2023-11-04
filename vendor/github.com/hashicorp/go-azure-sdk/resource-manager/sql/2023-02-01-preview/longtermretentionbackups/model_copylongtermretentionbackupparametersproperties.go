package longtermretentionbackups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CopyLongTermRetentionBackupParametersProperties struct {
	TargetBackupStorageRedundancy        *BackupStorageRedundancy `json:"targetBackupStorageRedundancy,omitempty"`
	TargetDatabaseName                   *string                  `json:"targetDatabaseName,omitempty"`
	TargetResourceGroup                  *string                  `json:"targetResourceGroup,omitempty"`
	TargetServerFullyQualifiedDomainName *string                  `json:"targetServerFullyQualifiedDomainName,omitempty"`
	TargetServerResourceId               *string                  `json:"targetServerResourceId,omitempty"`
	TargetSubscriptionId                 *string                  `json:"targetSubscriptionId,omitempty"`
}
