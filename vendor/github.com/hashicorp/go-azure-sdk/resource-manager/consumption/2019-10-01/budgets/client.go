package budgets

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BudgetsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewBudgetsClientWithBaseURI(endpoint string) BudgetsClient {
	return BudgetsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
