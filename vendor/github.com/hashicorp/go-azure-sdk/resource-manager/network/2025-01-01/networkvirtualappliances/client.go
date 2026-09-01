package networkvirtualappliances

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkVirtualAppliancesClient struct {
	Client *resourcemanager.Client
}

func NewNetworkVirtualAppliancesClientWithBaseURI(sdkApi sdkEnv.Api) (*NetworkVirtualAppliancesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "networkvirtualappliances", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating NetworkVirtualAppliancesClient: %+v", err)
	}

	return &NetworkVirtualAppliancesClient{
		Client: client,
	}, nil
}
