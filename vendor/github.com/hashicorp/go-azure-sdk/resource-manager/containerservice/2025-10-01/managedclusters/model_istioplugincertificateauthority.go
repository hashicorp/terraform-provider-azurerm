package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IstioPluginCertificateAuthority struct {
	CertChainObjectName *string `json:"certChainObjectName,omitempty"`
	CertObjectName      *string `json:"certObjectName,omitempty"`
	KeyObjectName       *string `json:"keyObjectName,omitempty"`
	KeyVaultId          *string `json:"keyVaultId,omitempty"`
	RootCertObjectName  *string `json:"rootCertObjectName,omitempty"`
}
