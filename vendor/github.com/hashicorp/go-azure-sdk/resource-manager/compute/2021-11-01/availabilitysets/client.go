package availabilitysets

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AvailabilitySetsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewAvailabilitySetsClientWithBaseURI(endpoint string) AvailabilitySetsClient {
	return AvailabilitySetsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
