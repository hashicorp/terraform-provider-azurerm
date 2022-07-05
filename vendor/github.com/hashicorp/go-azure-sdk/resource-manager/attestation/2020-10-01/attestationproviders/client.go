package attestationproviders

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AttestationProvidersClient struct {
	Client  autorest.Client
	baseUri string
}

func NewAttestationProvidersClientWithBaseURI(endpoint string) AttestationProvidersClient {
	return AttestationProvidersClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
