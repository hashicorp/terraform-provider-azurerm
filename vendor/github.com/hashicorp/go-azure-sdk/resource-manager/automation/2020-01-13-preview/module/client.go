package module

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ModuleClient struct {
	Client  autorest.Client
	baseUri string
}

func NewModuleClientWithBaseURI(endpoint string) ModuleClient {
	return ModuleClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
