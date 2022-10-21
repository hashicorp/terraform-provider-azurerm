package managementpolicies

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagementPoliciesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewManagementPoliciesClientWithBaseURI(endpoint string) ManagementPoliciesClient {
	return ManagementPoliciesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
