package networksecurityperimeterloggingconfigurations

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkSecurityPerimeterLoggingConfigurationsClient struct {
	Client *resourcemanager.Client
}

func NewNetworkSecurityPerimeterLoggingConfigurationsClientWithBaseURI(sdkApi sdkEnv.Api) (*NetworkSecurityPerimeterLoggingConfigurationsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "networksecurityperimeterloggingconfigurations", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating NetworkSecurityPerimeterLoggingConfigurationsClient: %+v", err)
	}

	return &NetworkSecurityPerimeterLoggingConfigurationsClient{
		Client: client,
	}, nil
}
