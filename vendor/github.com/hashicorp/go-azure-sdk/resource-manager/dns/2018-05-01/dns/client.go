package dns

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DNSClient struct {
	Client  autorest.Client
	baseUri string
}

func NewDNSClientWithBaseURI(endpoint string) DNSClient {
	return DNSClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
