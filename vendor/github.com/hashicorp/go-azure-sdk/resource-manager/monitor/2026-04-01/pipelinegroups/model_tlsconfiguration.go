package pipelinegroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TlsConfiguration struct {
	ClientCa       *CertificateSource  `json:"clientCa,omitempty"`
	Mode           *TlsMode            `json:"mode,omitempty"`
	Name           string              `json:"name"`
	TlsCertificate *CertificateWithKey `json:"tlsCertificate,omitempty"`
}
