package agentversions

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AgentVersionsClient struct {
	Client *resourcemanager.Client
}

func NewAgentVersionsClientWithBaseURI(sdkApi sdkEnv.Api) (*AgentVersionsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "agentversions", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AgentVersionsClient: %+v", err)
	}

	return &AgentVersionsClient{
		Client: client,
	}, nil
}
