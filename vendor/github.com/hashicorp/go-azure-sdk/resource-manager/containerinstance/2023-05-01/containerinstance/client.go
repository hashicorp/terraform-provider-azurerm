package containerinstance

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerInstanceClient struct {
	Client *resourcemanager.Client
}

func NewContainerInstanceClientWithBaseURI(sdkApi sdkEnv.Api) (*ContainerInstanceClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "containerinstance", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ContainerInstanceClient: %+v", err)
	}

	return &ContainerInstanceClient{
		Client: client,
	}, nil
}
