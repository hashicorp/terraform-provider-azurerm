package serverendpointresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FilesNotTieringError struct {
	ErrorCode *int64 `json:"errorCode,omitempty"`
	FileCount *int64 `json:"fileCount,omitempty"`
}
