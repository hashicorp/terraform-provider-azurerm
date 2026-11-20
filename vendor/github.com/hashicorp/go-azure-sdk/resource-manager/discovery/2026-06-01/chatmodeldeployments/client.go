package chatmodeldeployments

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ChatModelDeploymentsClient struct {
	Client *resourcemanager.Client
}

func NewChatModelDeploymentsClientWithBaseURI(sdkApi sdkEnv.Api) (*ChatModelDeploymentsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "chatmodeldeployments", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ChatModelDeploymentsClient: %+v", err)
	}

	return &ChatModelDeploymentsClient{
		Client: client,
	}, nil
}
