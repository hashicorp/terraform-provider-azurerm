package tableserviceproperties

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TableServicePropertiesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewTableServicePropertiesClientWithBaseURI(endpoint string) TableServicePropertiesClient {
	return TableServicePropertiesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
