package subgroups

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubgroupsClient struct {
	Client *resourcemanager.Client
}

func NewSubgroupsClientWithBaseURI(sdkApi sdkEnv.Api) (*SubgroupsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "subgroups", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SubgroupsClient: %+v", err)
	}

	return &SubgroupsClient{
		Client: client,
	}, nil
}
