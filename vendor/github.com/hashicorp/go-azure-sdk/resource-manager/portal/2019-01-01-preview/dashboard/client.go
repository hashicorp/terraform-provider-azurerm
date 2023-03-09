package dashboard

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DashboardClient struct {
	Client  autorest.Client
	baseUri string
}

func NewDashboardClientWithBaseURI(endpoint string) DashboardClient {
	return DashboardClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
