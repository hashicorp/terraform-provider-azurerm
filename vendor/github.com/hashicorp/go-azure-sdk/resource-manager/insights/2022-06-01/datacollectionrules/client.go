package datacollectionrules

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataCollectionRulesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewDataCollectionRulesClientWithBaseURI(endpoint string) DataCollectionRulesClient {
	return DataCollectionRulesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
