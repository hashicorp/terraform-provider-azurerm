package objectreplicationpolicies

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ObjectReplicationPoliciesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewObjectReplicationPoliciesClientWithBaseURI(endpoint string) ObjectReplicationPoliciesClient {
	return ObjectReplicationPoliciesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
