package monitors

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MonitorsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewMonitorsClientWithBaseURI(endpoint string) MonitorsClient {
	return MonitorsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
