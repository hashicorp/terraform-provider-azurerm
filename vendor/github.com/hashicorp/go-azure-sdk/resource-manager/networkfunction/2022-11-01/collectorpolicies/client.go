package collectorpolicies

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CollectorPoliciesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewCollectorPoliciesClientWithBaseURI(endpoint string) CollectorPoliciesClient {
	return CollectorPoliciesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
