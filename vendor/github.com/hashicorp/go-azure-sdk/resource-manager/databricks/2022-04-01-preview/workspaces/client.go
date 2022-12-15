package workspaces

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspacesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewWorkspacesClientWithBaseURI(endpoint string) WorkspacesClient {
	return WorkspacesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
