package deploymentscripts

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentScriptsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewDeploymentScriptsClientWithBaseURI(endpoint string) DeploymentScriptsClient {
	return DeploymentScriptsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
