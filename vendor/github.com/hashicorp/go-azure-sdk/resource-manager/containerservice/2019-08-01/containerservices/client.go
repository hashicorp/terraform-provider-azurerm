package containerservices

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerServicesClient struct {
	Client *resourcemanager.Client
}

func NewContainerServicesClientWithBaseURI(sdkApi sdkEnv.Api) (*ContainerServicesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "containerservices", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ContainerServicesClient: %+v", err)
	}

	return &ContainerServicesClient{
		Client: client,
	}, nil
}
