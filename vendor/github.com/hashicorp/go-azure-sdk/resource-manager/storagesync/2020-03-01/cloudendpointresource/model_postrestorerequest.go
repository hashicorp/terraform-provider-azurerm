package cloudendpointresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PostRestoreRequest struct {
	AzureFileShareUri       *string            `json:"azureFileShareUri,omitempty"`
	FailedFileList          *string            `json:"failedFileList,omitempty"`
	Partition               *string            `json:"partition,omitempty"`
	ReplicaGroup            *string            `json:"replicaGroup,omitempty"`
	RequestId               *string            `json:"requestId,omitempty"`
	RestoreFileSpec         *[]RestoreFileSpec `json:"restoreFileSpec,omitempty"`
	SourceAzureFileShareUri *string            `json:"sourceAzureFileShareUri,omitempty"`
	Status                  *string            `json:"status,omitempty"`
}
