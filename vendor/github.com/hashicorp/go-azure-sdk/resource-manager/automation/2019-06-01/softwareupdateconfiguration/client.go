package softwareupdateconfiguration

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SoftwareUpdateConfigurationClient struct {
	Client *resourcemanager.Client
}

func NewSoftwareUpdateConfigurationClientWithBaseURI(sdkApi sdkEnv.Api) (*SoftwareUpdateConfigurationClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "softwareupdateconfiguration", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SoftwareUpdateConfigurationClient: %+v", err)
	}

	return &SoftwareUpdateConfigurationClient{
		Client: client,
	}, nil
}
