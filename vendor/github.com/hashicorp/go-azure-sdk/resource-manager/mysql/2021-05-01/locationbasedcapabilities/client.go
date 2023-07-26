package locationbasedcapabilities

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LocationBasedCapabilitiesClient struct {
	Client *resourcemanager.Client
}

func NewLocationBasedCapabilitiesClientWithBaseURI(api environments.Api) (*LocationBasedCapabilitiesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "locationbasedcapabilities", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating LocationBasedCapabilitiesClient: %+v", err)
	}

	return &LocationBasedCapabilitiesClient{
		Client: client,
	}, nil
}
