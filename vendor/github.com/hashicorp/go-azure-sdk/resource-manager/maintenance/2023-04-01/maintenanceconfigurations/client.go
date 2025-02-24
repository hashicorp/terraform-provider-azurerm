package maintenanceconfigurations

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MaintenanceConfigurationsClient struct {
	Client *resourcemanager.Client
}

func NewMaintenanceConfigurationsClientWithBaseURI(sdkApi sdkEnv.Api) (*MaintenanceConfigurationsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "maintenanceconfigurations", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating MaintenanceConfigurationsClient: %+v", err)
	}

	return &MaintenanceConfigurationsClient{
		Client: client,
	}, nil
}
