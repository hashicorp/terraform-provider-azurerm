package attacheddatanetwork

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PortRange struct {
	MaxPort *int64 `json:"maxPort,omitempty"`
	MinPort *int64 `json:"minPort,omitempty"`
}
