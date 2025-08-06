package hubs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WnsCredentialProperties struct {
	CertificateKey      *string `json:"certificateKey,omitempty"`
	PackageSid          *string `json:"packageSid,omitempty"`
	SecretKey           *string `json:"secretKey,omitempty"`
	WindowsLiveEndpoint *string `json:"windowsLiveEndpoint,omitempty"`
	WnsCertificate      *string `json:"wnsCertificate,omitempty"`
}
