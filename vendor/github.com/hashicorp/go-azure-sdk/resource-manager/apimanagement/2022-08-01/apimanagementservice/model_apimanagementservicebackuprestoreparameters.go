package apimanagementservice

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiManagementServiceBackupRestoreParameters struct {
	AccessKey      *string     `json:"accessKey,omitempty"`
	AccessType     *AccessType `json:"accessType,omitempty"`
	BackupName     string      `json:"backupName"`
	ClientId       *string     `json:"clientId,omitempty"`
	ContainerName  string      `json:"containerName"`
	StorageAccount string      `json:"storageAccount"`
}
