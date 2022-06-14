package nameavailability

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NameAvailabilityClient struct {
	Client  autorest.Client
	baseUri string
}

func NewNameAvailabilityClientWithBaseURI(endpoint string) NameAvailabilityClient {
	return NameAvailabilityClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
