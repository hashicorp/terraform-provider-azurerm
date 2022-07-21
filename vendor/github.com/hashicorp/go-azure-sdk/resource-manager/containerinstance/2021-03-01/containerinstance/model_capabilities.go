package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Capabilities struct {
	Capabilities  *CapabilitiesCapabilities `json:"capabilities,omitempty"`
	Gpu           *string                   `json:"gpu,omitempty"`
	IpAddressType *string                   `json:"ipAddressType,omitempty"`
	Location      *string                   `json:"location,omitempty"`
	OsType        *string                   `json:"osType,omitempty"`
	ResourceType  *string                   `json:"resourceType,omitempty"`
}
