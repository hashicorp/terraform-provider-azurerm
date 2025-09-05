package managedcluster

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClientCertificate struct {
	CommonName       *string `json:"commonName,omitempty"`
	IsAdmin          bool    `json:"isAdmin"`
	IssuerThumbprint *string `json:"issuerThumbprint,omitempty"`
	Thumbprint       *string `json:"thumbprint,omitempty"`
}
