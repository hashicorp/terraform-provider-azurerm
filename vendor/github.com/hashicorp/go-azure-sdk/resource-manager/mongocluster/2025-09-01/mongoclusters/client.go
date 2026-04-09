package mongoclusters

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MongoClustersClient struct {
	Client *resourcemanager.Client
}

func NewMongoClustersClientWithBaseURI(sdkApi sdkEnv.Api) (*MongoClustersClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "mongoclusters", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating MongoClustersClient: %+v", err)
	}

	return &MongoClustersClient{
		Client: client,
	}, nil
}
