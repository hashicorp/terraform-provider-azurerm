package nodetype

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NodeTypeClient struct {
	Client *resourcemanager.Client
}

func NewNodeTypeClientWithBaseURI(sdkApi sdkEnv.Api) (*NodeTypeClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "nodetype", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating NodeTypeClient: %+v", err)
	}

	return &NodeTypeClient{
		Client: client,
	}, nil
}
