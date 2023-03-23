package automationaccount

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutomationAccountClient struct {
	Client  autorest.Client
	baseUri string
}

func NewAutomationAccountClientWithBaseURI(endpoint string) AutomationAccountClient {
	return AutomationAccountClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
