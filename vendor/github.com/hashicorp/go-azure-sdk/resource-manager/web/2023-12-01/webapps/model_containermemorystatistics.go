package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerMemoryStatistics struct {
	Limit    *int64 `json:"limit,omitempty"`
	MaxUsage *int64 `json:"maxUsage,omitempty"`
	Usage    *int64 `json:"usage,omitempty"`
}
