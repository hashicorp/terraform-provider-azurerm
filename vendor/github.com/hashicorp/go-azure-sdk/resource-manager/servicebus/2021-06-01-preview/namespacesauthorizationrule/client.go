package namespacesauthorizationrule

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NamespacesAuthorizationRuleClient struct {
	Client  autorest.Client
	baseUri string
}

func NewNamespacesAuthorizationRuleClientWithBaseURI(endpoint string) NamespacesAuthorizationRuleClient {
	return NamespacesAuthorizationRuleClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
