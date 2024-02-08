package eventsubscriptions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageBlobDeadLetterDestinationProperties struct {
	BlobContainerName *string `json:"blobContainerName,omitempty"`
	ResourceId        *string `json:"resourceId,omitempty"`
}
