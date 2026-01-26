package certificates

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificateBundle struct {
	Attributes        *CertificateAttributes `json:"attributes,omitempty"`
	Cer               *string                `json:"cer,omitempty"`
	ContentType       *string                `json:"contentType,omitempty"`
	Id                *string                `json:"id,omitempty"`
	Kid               *string                `json:"kid,omitempty"`
	Policy            *CertificatePolicy     `json:"policy,omitempty"`
	PreserveCertOrder *bool                  `json:"preserveCertOrder,omitempty"`
	Sid               *string                `json:"sid,omitempty"`
	Tags              *map[string]string     `json:"tags,omitempty"`
	X5t               *string                `json:"x5t,omitempty"`
}
