package forwardingrules

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ForwardingRulesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewForwardingRulesClientWithBaseURI(endpoint string) ForwardingRulesClient {
	return ForwardingRulesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
