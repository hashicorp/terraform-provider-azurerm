package configurations

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigurationsClient struct {
	Client *resourcemanager.Client
}

func NewConfigurationsClientWithBaseURI(sdkApi sdkEnv.Api) (*ConfigurationsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "configurations", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ConfigurationsClient: %+v", err)
	}

	return &ConfigurationsClient{
		Client: client,
	}, nil
}
