package intelligencepacks

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntelligencePacksClient struct {
	Client  autorest.Client
	baseUri string
}

func NewIntelligencePacksClientWithBaseURI(endpoint string) IntelligencePacksClient {
	return IntelligencePacksClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
