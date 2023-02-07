package datastore

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatastoreClient struct {
	Client  autorest.Client
	baseUri string
}

func NewDatastoreClientWithBaseURI(endpoint string) DatastoreClient {
	return DatastoreClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
