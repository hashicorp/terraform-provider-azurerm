package automationrules

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutomationRulesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewAutomationRulesClientWithBaseURI(endpoint string) AutomationRulesClient {
	return AutomationRulesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
