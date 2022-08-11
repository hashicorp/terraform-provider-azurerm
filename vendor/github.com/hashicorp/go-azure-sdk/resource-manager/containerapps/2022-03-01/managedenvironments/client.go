package managedenvironments

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedEnvironmentsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewManagedEnvironmentsClientWithBaseURI(endpoint string) ManagedEnvironmentsClient {
	return ManagedEnvironmentsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
