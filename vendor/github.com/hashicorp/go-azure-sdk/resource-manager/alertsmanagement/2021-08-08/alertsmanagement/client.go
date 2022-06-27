package alertsmanagement

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlertsManagementClient struct {
	Client  autorest.Client
	baseUri string
}

func NewAlertsManagementClientWithBaseURI(endpoint string) AlertsManagementClient {
	return AlertsManagementClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
