package cloudendpointresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PreRestoreRequest struct {
	AzureFileShareUri                        *string            `json:"azureFileShareUri,omitempty"`
	BackupMetadataPropertyBag                *string            `json:"backupMetadataPropertyBag,omitempty"`
	Partition                                *string            `json:"partition,omitempty"`
	PauseWaitForSyncDrainTimePeriodInSeconds *int64             `json:"pauseWaitForSyncDrainTimePeriodInSeconds,omitempty"`
	ReplicaGroup                             *string            `json:"replicaGroup,omitempty"`
	RequestId                                *string            `json:"requestId,omitempty"`
	RestoreFileSpec                          *[]RestoreFileSpec `json:"restoreFileSpec,omitempty"`
	SourceAzureFileShareUri                  *string            `json:"sourceAzureFileShareUri,omitempty"`
	Status                                   *string            `json:"status,omitempty"`
}
