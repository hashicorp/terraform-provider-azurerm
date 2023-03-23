package managementlocks

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagementLocksClient struct {
	Client  autorest.Client
	baseUri string
}

func NewManagementLocksClientWithBaseURI(endpoint string) ManagementLocksClient {
	return ManagementLocksClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
