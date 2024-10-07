package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateAccessVirtualNetwork struct {
	Key        *int64                 `json:"key,omitempty"`
	Name       *string                `json:"name,omitempty"`
	ResourceId *string                `json:"resourceId,omitempty"`
	Subnets    *[]PrivateAccessSubnet `json:"subnets,omitempty"`
}
