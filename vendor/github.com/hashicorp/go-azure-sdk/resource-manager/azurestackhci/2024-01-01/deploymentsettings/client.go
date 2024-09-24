package deploymentsettings

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentSettingsClient struct {
	Client *resourcemanager.Client
}

func NewDeploymentSettingsClientWithBaseURI(sdkApi sdkEnv.Api) (*DeploymentSettingsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "deploymentsettings", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DeploymentSettingsClient: %+v", err)
	}

	return &DeploymentSettingsClient{
		Client: client,
	}, nil
}
