package servicelinker

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceLinkerClient struct {
	Client *resourcemanager.Client
}

func NewServiceLinkerClientWithBaseURI(sdkApi sdkEnv.Api) (*ServiceLinkerClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "servicelinker", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ServiceLinkerClient: %+v", err)
	}

	return &ServiceLinkerClient{
		Client: client,
	}, nil
}
