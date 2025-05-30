package solution

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SolutionClient struct {
	Client *resourcemanager.Client
}

func NewSolutionClientWithBaseURI(sdkApi sdkEnv.Api) (*SolutionClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "solution", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SolutionClient: %+v", err)
	}

	return &SolutionClient{
		Client: client,
	}, nil
}
