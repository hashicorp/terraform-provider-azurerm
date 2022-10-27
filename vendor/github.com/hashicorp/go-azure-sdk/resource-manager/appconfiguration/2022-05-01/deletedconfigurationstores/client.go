package deletedconfigurationstores

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeletedConfigurationStoresClient struct {
	Client  autorest.Client
	baseUri string
}

func NewDeletedConfigurationStoresClientWithBaseURI(endpoint string) DeletedConfigurationStoresClient {
	return DeletedConfigurationStoresClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
