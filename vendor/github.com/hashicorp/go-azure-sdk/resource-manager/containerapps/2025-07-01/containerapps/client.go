package containerapps

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerAppsClient struct {
	Client *resourcemanager.Client
}

func NewContainerAppsClientWithBaseURI(sdkApi sdkEnv.Api) (*ContainerAppsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "containerapps", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ContainerAppsClient: %+v", err)
	}

	return &ContainerAppsClient{
		Client: client,
	}, nil
}
