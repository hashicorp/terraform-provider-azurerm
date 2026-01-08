package broker

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackendChain struct {
	Partitions       int64  `json:"partitions"`
	RedundancyFactor int64  `json:"redundancyFactor"`
	Workers          *int64 `json:"workers,omitempty"`
}
