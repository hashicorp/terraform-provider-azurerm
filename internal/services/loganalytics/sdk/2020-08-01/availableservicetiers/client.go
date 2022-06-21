package availableservicetiers

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AvailableServiceTiersClient struct {
	Client  autorest.Client
	baseUri string
}

func NewAvailableServiceTiersClientWithBaseURI(endpoint string) AvailableServiceTiersClient {
	return AvailableServiceTiersClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
