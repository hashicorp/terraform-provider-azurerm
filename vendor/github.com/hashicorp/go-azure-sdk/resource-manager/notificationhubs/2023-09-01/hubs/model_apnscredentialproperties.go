package hubs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApnsCredentialProperties struct {
	ApnsCertificate *string `json:"apnsCertificate,omitempty"`
	AppId           *string `json:"appId,omitempty"`
	AppName         *string `json:"appName,omitempty"`
	CertificateKey  *string `json:"certificateKey,omitempty"`
	Endpoint        string  `json:"endpoint"`
	KeyId           *string `json:"keyId,omitempty"`
	Thumbprint      *string `json:"thumbprint,omitempty"`
	Token           *string `json:"token,omitempty"`
}
