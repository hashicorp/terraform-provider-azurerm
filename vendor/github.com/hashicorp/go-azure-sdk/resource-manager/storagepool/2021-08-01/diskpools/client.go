package diskpools

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiskPoolsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewDiskPoolsClientWithBaseURI(endpoint string) DiskPoolsClient {
	return DiskPoolsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
