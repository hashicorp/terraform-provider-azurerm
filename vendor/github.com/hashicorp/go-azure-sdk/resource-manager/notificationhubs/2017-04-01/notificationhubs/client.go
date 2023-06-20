package notificationhubs

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NotificationHubsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewNotificationHubsClientWithBaseURI(endpoint string) NotificationHubsClient {
	return NotificationHubsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
