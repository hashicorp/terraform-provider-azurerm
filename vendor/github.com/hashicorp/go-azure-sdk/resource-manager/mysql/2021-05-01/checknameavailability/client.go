package checknameavailability

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CheckNameAvailabilityClient struct {
	Client  autorest.Client
	baseUri string
}

func NewCheckNameAvailabilityClientWithBaseURI(endpoint string) CheckNameAvailabilityClient {
	return CheckNameAvailabilityClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
