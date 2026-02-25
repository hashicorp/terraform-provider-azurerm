package exadbvmclusters

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExadbVMClustersClient struct {
	Client *resourcemanager.Client
}

func NewExadbVMClustersClientWithBaseURI(sdkApi sdkEnv.Api) (*ExadbVMClustersClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "exadbvmclusters", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ExadbVMClustersClient: %+v", err)
	}

	return &ExadbVMClustersClient{
		Client: client,
	}, nil
}
