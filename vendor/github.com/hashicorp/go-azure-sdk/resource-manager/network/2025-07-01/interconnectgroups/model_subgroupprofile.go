package interconnectgroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubgroupProfile struct {
	Scope  *SubgroupProfileScope `json:"scope,omitempty"`
	Size   *int64                `json:"size,omitempty"`
	VMSize string                `json:"vmSize"`
}
