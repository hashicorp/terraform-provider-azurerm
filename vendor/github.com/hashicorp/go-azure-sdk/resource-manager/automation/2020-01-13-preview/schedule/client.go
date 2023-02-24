package schedule

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScheduleClient struct {
	Client  autorest.Client
	baseUri string
}

func NewScheduleClientWithBaseURI(endpoint string) ScheduleClient {
	return ScheduleClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
