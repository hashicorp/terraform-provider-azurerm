package virtualwans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualWanSecurityProvider struct {
	Name *string                         `json:"name,omitempty"`
	Type *VirtualWanSecurityProviderType `json:"type,omitempty"`
	Url  *string                         `json:"url,omitempty"`
}
