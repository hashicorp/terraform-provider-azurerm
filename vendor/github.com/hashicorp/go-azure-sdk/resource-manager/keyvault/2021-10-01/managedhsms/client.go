package managedhsms

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedHsmsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewManagedHsmsClientWithBaseURI(endpoint string) ManagedHsmsClient {
	return ManagedHsmsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
