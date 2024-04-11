package publicmaintenanceconfigurations

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PublicMaintenanceConfigurationsClient struct {
	Client *resourcemanager.Client
}

func NewPublicMaintenanceConfigurationsClientWithBaseURI(sdkApi sdkEnv.Api) (*PublicMaintenanceConfigurationsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "publicmaintenanceconfigurations", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PublicMaintenanceConfigurationsClient: %+v", err)
	}

	return &PublicMaintenanceConfigurationsClient{
		Client: client,
	}, nil
}
