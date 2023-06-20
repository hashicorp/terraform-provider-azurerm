package alertrules

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlertRulesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewAlertRulesClientWithBaseURI(endpoint string) AlertRulesClient {
	return AlertRulesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
