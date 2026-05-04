package serverendpointresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerEndpointRecallError struct {
	Count     *int64 `json:"count,omitempty"`
	ErrorCode *int64 `json:"errorCode,omitempty"`
}
