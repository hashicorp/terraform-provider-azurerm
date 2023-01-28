package schedules

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SchedulesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewSchedulesClientWithBaseURI(endpoint string) SchedulesClient {
	return SchedulesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
