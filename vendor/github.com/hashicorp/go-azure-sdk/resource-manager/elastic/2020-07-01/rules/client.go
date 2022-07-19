package rules

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RulesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewRulesClientWithBaseURI(endpoint string) RulesClient {
	return RulesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
