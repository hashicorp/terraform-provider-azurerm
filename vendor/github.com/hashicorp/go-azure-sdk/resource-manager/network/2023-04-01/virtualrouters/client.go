package virtualrouters

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualRoutersClient struct {
	Client *resourcemanager.Client
}

func NewVirtualRoutersClientWithBaseURI(sdkApi sdkEnv.Api) (*VirtualRoutersClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "virtualrouters", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating VirtualRoutersClient: %+v", err)
	}

	return &VirtualRoutersClient{
		Client: client,
	}, nil
}
