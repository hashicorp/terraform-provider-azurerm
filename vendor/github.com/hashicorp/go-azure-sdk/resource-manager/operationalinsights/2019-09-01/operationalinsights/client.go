package operationalinsights

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OperationalInsightsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewOperationalInsightsClientWithBaseURI(endpoint string) OperationalInsightsClient {
	return OperationalInsightsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
