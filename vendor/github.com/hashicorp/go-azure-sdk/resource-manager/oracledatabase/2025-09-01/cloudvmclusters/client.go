package cloudvmclusters

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CloudVMClustersClient struct {
	Client *resourcemanager.Client
}

func NewCloudVMClustersClientWithBaseURI(sdkApi sdkEnv.Api) (*CloudVMClustersClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "cloudvmclusters", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating CloudVMClustersClient: %+v", err)
	}

	return &CloudVMClustersClient{
		Client: client,
	}, nil
}
