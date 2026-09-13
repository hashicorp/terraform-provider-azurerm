package certificates

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificateUpdateParameters struct {
	Attributes *CertificateAttributes `json:"attributes,omitempty"`
	Policy     *CertificatePolicy     `json:"policy,omitempty"`
	Tags       *map[string]string     `json:"tags,omitempty"`
}
