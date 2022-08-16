package workbooksapis

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkbooksAPIsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewWorkbooksAPIsClientWithBaseURI(endpoint string) WorkbooksAPIsClient {
	return WorkbooksAPIsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
