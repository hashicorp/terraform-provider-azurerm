package configurationprofiles

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigurationProfilesClient struct {
	Client *resourcemanager.Client
}

func NewConfigurationProfilesClientWithBaseURI(sdkApi sdkEnv.Api) (*ConfigurationProfilesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "configurationprofiles", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ConfigurationProfilesClient: %+v", err)
	}

	return &ConfigurationProfilesClient{
		Client: client,
	}, nil
}
