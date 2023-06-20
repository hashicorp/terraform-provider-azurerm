package lab

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LabClient struct {
	Client  autorest.Client
	baseUri string
}

func NewLabClientWithBaseURI(endpoint string) LabClient {
	return LabClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
