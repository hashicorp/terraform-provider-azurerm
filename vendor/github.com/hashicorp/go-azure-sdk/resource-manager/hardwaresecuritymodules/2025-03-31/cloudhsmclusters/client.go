package cloudhsmclusters

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CloudHsmClustersClient struct {
	Client *resourcemanager.Client
}

func NewCloudHsmClustersClientWithBaseURI(sdkApi sdkEnv.Api) (*CloudHsmClustersClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "cloudhsmclusters", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating CloudHsmClustersClient: %+v", err)
	}

	return &CloudHsmClustersClient{
		Client: client,
	}, nil
}
