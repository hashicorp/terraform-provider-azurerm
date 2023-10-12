package integrationserviceenvironments

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntegrationServiceEnvironmentsClient struct {
	Client *resourcemanager.Client
}

func NewIntegrationServiceEnvironmentsClientWithBaseURI(sdkApi sdkEnv.Api) (*IntegrationServiceEnvironmentsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "integrationserviceenvironments", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating IntegrationServiceEnvironmentsClient: %+v", err)
	}

	return &IntegrationServiceEnvironmentsClient{
		Client: client,
	}, nil
}
