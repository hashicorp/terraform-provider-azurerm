package availabilitygrouplisteners

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AvailabilityGroupListenersClient struct {
	Client  autorest.Client
	baseUri string
}

func NewAvailabilityGroupListenersClientWithBaseURI(endpoint string) AvailabilityGroupListenersClient {
	return AvailabilityGroupListenersClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
