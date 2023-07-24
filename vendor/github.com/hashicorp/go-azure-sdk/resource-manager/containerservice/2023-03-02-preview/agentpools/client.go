package agentpools

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AgentPoolsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewAgentPoolsClientWithBaseURI(endpoint string) AgentPoolsClient {
	return AgentPoolsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
