package topicsauthorizationrule

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TopicsAuthorizationRuleClient struct {
	Client  autorest.Client
	baseUri string
}

func NewTopicsAuthorizationRuleClientWithBaseURI(endpoint string) TopicsAuthorizationRuleClient {
	return TopicsAuthorizationRuleClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
