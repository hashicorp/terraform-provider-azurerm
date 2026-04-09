package resourceproviders

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualNetworkProfile struct {
	Id     string  `json:"id"`
	Name   *string `json:"name,omitempty"`
	Subnet *string `json:"subnet,omitempty"`
	Type   *string `json:"type,omitempty"`
}
