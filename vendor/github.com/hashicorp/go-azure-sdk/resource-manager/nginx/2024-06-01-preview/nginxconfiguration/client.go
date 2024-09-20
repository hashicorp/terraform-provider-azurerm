package nginxconfiguration

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NginxConfigurationClient struct {
	Client *resourcemanager.Client
}

func NewNginxConfigurationClientWithBaseURI(sdkApi sdkEnv.Api) (*NginxConfigurationClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "nginxconfiguration", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating NginxConfigurationClient: %+v", err)
	}

	return &NginxConfigurationClient{
		Client: client,
	}, nil
}
