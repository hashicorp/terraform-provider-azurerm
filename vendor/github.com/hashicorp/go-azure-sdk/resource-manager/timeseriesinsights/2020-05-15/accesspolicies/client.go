package accesspolicies

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccessPoliciesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewAccessPoliciesClientWithBaseURI(endpoint string) AccessPoliciesClient {
	return AccessPoliciesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
