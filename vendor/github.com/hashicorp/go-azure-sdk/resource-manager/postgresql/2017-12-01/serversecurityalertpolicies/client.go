package serversecurityalertpolicies

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerSecurityAlertPoliciesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewServerSecurityAlertPoliciesClientWithBaseURI(endpoint string) ServerSecurityAlertPoliciesClient {
	return ServerSecurityAlertPoliciesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
