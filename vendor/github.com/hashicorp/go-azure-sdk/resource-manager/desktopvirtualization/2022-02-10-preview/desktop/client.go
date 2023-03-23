package desktop

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DesktopClient struct {
	Client  autorest.Client
	baseUri string
}

func NewDesktopClientWithBaseURI(endpoint string) DesktopClient {
	return DesktopClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
