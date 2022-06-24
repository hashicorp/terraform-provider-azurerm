package savedsearches

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SavedSearchesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewSavedSearchesClientWithBaseURI(endpoint string) SavedSearchesClient {
	return SavedSearchesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
