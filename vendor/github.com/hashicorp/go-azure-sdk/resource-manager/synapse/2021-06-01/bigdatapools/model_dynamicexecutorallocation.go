package bigdatapools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DynamicExecutorAllocation struct {
	Enabled      *bool  `json:"enabled,omitempty"`
	MaxExecutors *int64 `json:"maxExecutors,omitempty"`
	MinExecutors *int64 `json:"minExecutors,omitempty"`
}
