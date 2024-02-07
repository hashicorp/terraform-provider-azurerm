package certificates

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificatesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewCertificatesClientWithBaseURI(endpoint string) CertificatesClient {
	return CertificatesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
