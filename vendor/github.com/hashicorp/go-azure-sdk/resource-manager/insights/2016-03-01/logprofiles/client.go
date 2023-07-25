package logprofiles

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LogProfilesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewLogProfilesClientWithBaseURI(endpoint string) LogProfilesClient {
	return LogProfilesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
