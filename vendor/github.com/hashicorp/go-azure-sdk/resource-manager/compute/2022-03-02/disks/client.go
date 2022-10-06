package disks

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DisksClient struct {
	Client  autorest.Client
	baseUri string
}

func NewDisksClientWithBaseURI(endpoint string) DisksClient {
	return DisksClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
