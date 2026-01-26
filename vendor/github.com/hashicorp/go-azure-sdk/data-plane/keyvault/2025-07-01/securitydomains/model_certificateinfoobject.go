package securitydomains

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificateInfoObject struct {
	Certificates []SecurityDomainJsonWebKey `json:"certificates"`
	Required     *int64                     `json:"required,omitempty"`
}
