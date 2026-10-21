package jobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OutputFileBlobContainerDestination struct {
	ContainerURL      string                        `json:"containerUrl"`
	IdentityReference *ComputeNodeIdentityReference `json:"identityReference,omitempty"`
	Path              *string                       `json:"path,omitempty"`
	UploadHeaders     *[]HTTPHeader                 `json:"uploadHeaders,omitempty"`
}
