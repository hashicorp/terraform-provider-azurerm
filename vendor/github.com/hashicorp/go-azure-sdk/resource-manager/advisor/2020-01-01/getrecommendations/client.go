package getrecommendations

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetRecommendationsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewGetRecommendationsClientWithBaseURI(endpoint string) GetRecommendationsClient {
	return GetRecommendationsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
