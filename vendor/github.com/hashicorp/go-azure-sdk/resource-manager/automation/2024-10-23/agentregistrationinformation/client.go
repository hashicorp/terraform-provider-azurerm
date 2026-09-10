package agentregistrationinformation

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AgentRegistrationInformationClient struct {
	Client *resourcemanager.Client
}

func NewAgentRegistrationInformationClientWithBaseURI(sdkApi sdkEnv.Api) (*AgentRegistrationInformationClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "agentregistrationinformation", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AgentRegistrationInformationClient: %+v", err)
	}

	return &AgentRegistrationInformationClient{
		Client: client,
	}, nil
}
