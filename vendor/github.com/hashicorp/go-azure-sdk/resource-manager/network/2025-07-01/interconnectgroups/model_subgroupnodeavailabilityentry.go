package interconnectgroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubgroupNodeAvailabilityEntry struct {
	Count              *int64  `json:"count,omitempty"`
	InServiceNodeCount *int64  `json:"inServiceNodeCount,omitempty"`
	InUseNodeCount     *int64  `json:"inUseNodeCount,omitempty"`
	InternalSubgroupId *string `json:"internalSubgroupId,omitempty"`
	Name               *string `json:"name,omitempty"`
}
