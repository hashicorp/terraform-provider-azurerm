package pool

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureBlobFileSystemConfiguration struct {
	AccountKey        *string                       `json:"accountKey,omitempty"`
	AccountName       string                        `json:"accountName"`
	BlobfuseOptions   *string                       `json:"blobfuseOptions,omitempty"`
	ContainerName     string                        `json:"containerName"`
	IdentityReference *ComputeNodeIdentityReference `json:"identityReference,omitempty"`
	RelativeMountPath string                        `json:"relativeMountPath"`
	SasKey            *string                       `json:"sasKey,omitempty"`
}
