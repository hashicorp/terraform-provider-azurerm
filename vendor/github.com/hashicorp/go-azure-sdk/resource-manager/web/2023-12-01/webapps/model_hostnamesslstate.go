package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HostNameSslState struct {
	HostType   *HostType `json:"hostType,omitempty"`
	Name       *string   `json:"name,omitempty"`
	SslState   *SslState `json:"sslState,omitempty"`
	Thumbprint *string   `json:"thumbprint,omitempty"`
	ToUpdate   *bool     `json:"toUpdate,omitempty"`
	VirtualIP  *string   `json:"virtualIP,omitempty"`
}
