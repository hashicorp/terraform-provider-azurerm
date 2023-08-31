package configurationsupdate

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigurationsUpdateClient struct {
	Client *resourcemanager.Client
}

func NewConfigurationsUpdateClientWithBaseURI(sdkApi sdkEnv.Api) (*ConfigurationsUpdateClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "configurationsupdate", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ConfigurationsUpdateClient: %+v", err)
	}

	return &ConfigurationsUpdateClient{
		Client: client,
	}, nil
}
