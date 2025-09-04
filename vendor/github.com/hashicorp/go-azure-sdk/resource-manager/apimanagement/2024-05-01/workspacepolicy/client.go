package workspacepolicy

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspacePolicyClient struct {
	Client *resourcemanager.Client
}

func NewWorkspacePolicyClientWithBaseURI(sdkApi sdkEnv.Api) (*WorkspacePolicyClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "workspacepolicy", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating WorkspacePolicyClient: %+v", err)
	}

	return &WorkspacePolicyClient{
		Client: client,
	}, nil
}
