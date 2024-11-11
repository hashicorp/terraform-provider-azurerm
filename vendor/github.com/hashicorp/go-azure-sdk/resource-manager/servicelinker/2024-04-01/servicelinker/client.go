package servicelinker

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServicelinkerClient struct {
	Client *resourcemanager.Client
}

func NewServicelinkerClientWithBaseURI(sdkApi sdkEnv.Api) (*ServicelinkerClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "servicelinker", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ServicelinkerClient: %+v", err)
	}

	return &ServicelinkerClient{
		Client: client,
	}, nil
}
