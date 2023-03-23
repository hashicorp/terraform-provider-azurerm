package containerservices

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerServicesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewContainerServicesClientWithBaseURI(endpoint string) ContainerServicesClient {
	return ContainerServicesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
