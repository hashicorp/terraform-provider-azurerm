package locationbasedcapabilities

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LocationBasedCapabilitiesClient struct {
	Client *resourcemanager.Client
}

func NewLocationBasedCapabilitiesClientWithBaseURI(sdkApi sdkEnv.Api) (*LocationBasedCapabilitiesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "locationbasedcapabilities", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating LocationBasedCapabilitiesClient: %+v", err)
	}

	return &LocationBasedCapabilitiesClient{
		Client: client,
	}, nil
}
