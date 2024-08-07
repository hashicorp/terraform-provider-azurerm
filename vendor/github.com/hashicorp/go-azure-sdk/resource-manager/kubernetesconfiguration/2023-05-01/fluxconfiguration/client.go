package fluxconfiguration

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FluxConfigurationClient struct {
	Client *resourcemanager.Client
}

func NewFluxConfigurationClientWithBaseURI(sdkApi sdkEnv.Api) (*FluxConfigurationClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "fluxconfiguration", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating FluxConfigurationClient: %+v", err)
	}

	return &FluxConfigurationClient{
		Client: client,
	}, nil
}
