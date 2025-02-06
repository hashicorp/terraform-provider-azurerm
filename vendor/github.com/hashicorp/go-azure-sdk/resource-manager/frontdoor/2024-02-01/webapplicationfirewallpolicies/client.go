package webapplicationfirewallpolicies

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebApplicationFirewallPoliciesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewWebApplicationFirewallPoliciesClientWithBaseURI(endpoint string) WebApplicationFirewallPoliciesClient {
	return WebApplicationFirewallPoliciesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
