package containerappsrevisions

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerAppsRevisionsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewContainerAppsRevisionsClientWithBaseURI(endpoint string) ContainerAppsRevisionsClient {
	return ContainerAppsRevisionsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
