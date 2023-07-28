package onlinedeployment

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OnlineRequestSettings struct {
	MaxConcurrentRequestsPerInstance *int64  `json:"maxConcurrentRequestsPerInstance,omitempty"`
	MaxQueueWait                     *string `json:"maxQueueWait,omitempty"`
	RequestTimeout                   *string `json:"requestTimeout,omitempty"`
}
