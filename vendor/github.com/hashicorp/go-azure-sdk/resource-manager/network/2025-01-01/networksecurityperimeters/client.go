package networksecurityperimeters

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkSecurityPerimetersClient struct {
	Client *resourcemanager.Client
}

func NewNetworkSecurityPerimetersClientWithBaseURI(sdkApi sdkEnv.Api) (*NetworkSecurityPerimetersClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "networksecurityperimeters", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating NetworkSecurityPerimetersClient: %+v", err)
	}

	return &NetworkSecurityPerimetersClient{
		Client: client,
	}, nil
}
