package nginxdeployment

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NginxDeploymentClient struct {
	Client  autorest.Client
	baseUri string
}

func NewNginxDeploymentClientWithBaseURI(endpoint string) NginxDeploymentClient {
	return NginxDeploymentClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
