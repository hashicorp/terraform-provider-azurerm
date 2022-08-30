package remediations

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RemediationsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewRemediationsClientWithBaseURI(endpoint string) RemediationsClient {
	return RemediationsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
