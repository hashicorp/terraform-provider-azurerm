package containerinstance

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerInstanceClient struct {
	Client  autorest.Client
	baseUri string
}

func NewContainerInstanceClientWithBaseURI(endpoint string) ContainerInstanceClient {
	return ContainerInstanceClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
