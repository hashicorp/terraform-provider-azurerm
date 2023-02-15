package jobschedule

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobScheduleClient struct {
	Client  autorest.Client
	baseUri string
}

func NewJobScheduleClientWithBaseURI(endpoint string) JobScheduleClient {
	return JobScheduleClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
