package schemaregistry

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SchemaRegistryClient struct {
	Client  autorest.Client
	baseUri string
}

func NewSchemaRegistryClientWithBaseURI(endpoint string) SchemaRegistryClient {
	return SchemaRegistryClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
