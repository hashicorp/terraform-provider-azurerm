package tagrules

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TagRulesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewTagRulesClientWithBaseURI(endpoint string) TagRulesClient {
	return TagRulesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
