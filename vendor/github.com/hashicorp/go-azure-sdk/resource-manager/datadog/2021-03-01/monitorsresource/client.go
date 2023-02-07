package monitorsresource

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MonitorsResourceClient struct {
	Client  autorest.Client
	baseUri string
}

func NewMonitorsResourceClientWithBaseURI(endpoint string) MonitorsResourceClient {
	return MonitorsResourceClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
