package protecteditems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DistributedNodesInfo struct {
	ErrorDetail      *ErrorDetail `json:"errorDetail,omitempty"`
	NodeName         *string      `json:"nodeName,omitempty"`
	SourceResourceId *string      `json:"sourceResourceId,omitempty"`
	Status           *string      `json:"status,omitempty"`
}
