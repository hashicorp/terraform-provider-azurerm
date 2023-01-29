package deployments

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewDeploymentsClientWithBaseURI(endpoint string) DeploymentsClient {
	return DeploymentsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
