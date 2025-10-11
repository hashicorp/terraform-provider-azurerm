package networksecurityperimeterprofiles

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkSecurityPerimeterProfilesClient struct {
	Client *resourcemanager.Client
}

func NewNetworkSecurityPerimeterProfilesClientWithBaseURI(sdkApi sdkEnv.Api) (*NetworkSecurityPerimeterProfilesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "networksecurityperimeterprofiles", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating NetworkSecurityPerimeterProfilesClient: %+v", err)
	}

	return &NetworkSecurityPerimeterProfilesClient{
		Client: client,
	}, nil
}
