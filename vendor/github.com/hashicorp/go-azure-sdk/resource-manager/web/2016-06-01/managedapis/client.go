package managedapis

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedAPIsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewManagedAPIsClientWithBaseURI(endpoint string) ManagedAPIsClient {
	return ManagedAPIsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
