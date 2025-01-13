package networksecurityperimeterconfiguration

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkSecurityPerimeterConfigurationClient struct {
	Client *resourcemanager.Client
}

func NewNetworkSecurityPerimeterConfigurationClientWithBaseURI(sdkApi sdkEnv.Api) (*NetworkSecurityPerimeterConfigurationClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "networksecurityperimeterconfiguration", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating NetworkSecurityPerimeterConfigurationClient: %+v", err)
	}

	return &NetworkSecurityPerimeterConfigurationClient{
		Client: client,
	}, nil
}
