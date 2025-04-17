package nodetype

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NodeTypeNatConfig struct {
	BackendPort            *int64 `json:"backendPort,omitempty"`
	FrontendPortRangeEnd   *int64 `json:"frontendPortRangeEnd,omitempty"`
	FrontendPortRangeStart *int64 `json:"frontendPortRangeStart,omitempty"`
}
