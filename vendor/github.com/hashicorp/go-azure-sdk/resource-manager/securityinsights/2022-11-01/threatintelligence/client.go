package threatintelligence

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ThreatIntelligenceClient struct {
	Client  autorest.Client
	baseUri string
}

func NewThreatIntelligenceClientWithBaseURI(endpoint string) ThreatIntelligenceClient {
	return ThreatIntelligenceClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
