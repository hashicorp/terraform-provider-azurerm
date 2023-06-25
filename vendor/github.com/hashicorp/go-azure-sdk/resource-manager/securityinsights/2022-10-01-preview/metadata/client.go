package metadata

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MetadataClient struct {
	Client  autorest.Client
	baseUri string
}

func NewMetadataClientWithBaseURI(endpoint string) MetadataClient {
	return MetadataClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
