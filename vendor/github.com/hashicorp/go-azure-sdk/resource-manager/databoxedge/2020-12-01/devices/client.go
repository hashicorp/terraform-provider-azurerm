package devices

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DevicesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewDevicesClientWithBaseURI(endpoint string) DevicesClient {
	return DevicesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
