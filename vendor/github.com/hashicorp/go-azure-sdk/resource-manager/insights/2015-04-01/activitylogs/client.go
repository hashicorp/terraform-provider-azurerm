package activitylogs

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActivityLogsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewActivityLogsClientWithBaseURI(endpoint string) ActivityLogsClient {
	return ActivityLogsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
