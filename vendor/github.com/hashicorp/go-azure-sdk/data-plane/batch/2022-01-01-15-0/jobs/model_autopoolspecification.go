package jobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoPoolSpecification struct {
	AutoPoolIdPrefix   *string            `json:"autoPoolIdPrefix,omitempty"`
	KeepAlive          *bool              `json:"keepAlive,omitempty"`
	Pool               *PoolSpecification `json:"pool,omitempty"`
	PoolLifetimeOption PoolLifetimeOption `json:"poolLifetimeOption"`
}
