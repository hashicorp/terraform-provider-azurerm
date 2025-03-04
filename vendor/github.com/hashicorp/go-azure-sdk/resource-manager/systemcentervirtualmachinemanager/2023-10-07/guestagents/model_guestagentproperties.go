package guestagents

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GuestAgentProperties struct {
	Credentials        *GuestCredential        `json:"credentials,omitempty"`
	CustomResourceName *string                 `json:"customResourceName,omitempty"`
	HTTPProxyConfig    *HTTPProxyConfiguration `json:"httpProxyConfig,omitempty"`
	ProvisioningAction *ProvisioningAction     `json:"provisioningAction,omitempty"`
	ProvisioningState  *ProvisioningState      `json:"provisioningState,omitempty"`
	Status             *string                 `json:"status,omitempty"`
	Uuid               *string                 `json:"uuid,omitempty"`
}
