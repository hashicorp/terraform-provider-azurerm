package serverendpointresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerEndpointFilesNotSyncingError struct {
	ErrorCode       *int64 `json:"errorCode,omitempty"`
	PersistentCount *int64 `json:"persistentCount,omitempty"`
	TransientCount  *int64 `json:"transientCount,omitempty"`
}
