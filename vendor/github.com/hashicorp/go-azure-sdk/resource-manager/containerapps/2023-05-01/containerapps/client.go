package containerapps

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerAppsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewContainerAppsClientWithBaseURI(endpoint string) ContainerAppsClient {
	return ContainerAppsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
