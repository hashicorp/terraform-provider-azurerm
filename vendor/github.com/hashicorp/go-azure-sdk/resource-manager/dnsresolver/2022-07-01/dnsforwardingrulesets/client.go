package dnsforwardingrulesets

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DnsForwardingRulesetsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewDnsForwardingRulesetsClientWithBaseURI(endpoint string) DnsForwardingRulesetsClient {
	return DnsForwardingRulesetsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
