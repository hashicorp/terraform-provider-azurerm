package dpscertificate

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DpsCertificateClient struct {
	Client  autorest.Client
	baseUri string
}

func NewDpsCertificateClientWithBaseURI(endpoint string) DpsCertificateClient {
	return DpsCertificateClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
