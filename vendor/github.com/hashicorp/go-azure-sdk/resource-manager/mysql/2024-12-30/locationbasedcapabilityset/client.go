package locationbasedcapabilityset

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LocationBasedCapabilitySetClient struct {
	Client *resourcemanager.Client
}

func NewLocationBasedCapabilitySetClientWithBaseURI(sdkApi sdkEnv.Api) (*LocationBasedCapabilitySetClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "locationbasedcapabilityset", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating LocationBasedCapabilitySetClient: %+v", err)
	}

	return &LocationBasedCapabilitySetClient{
		Client: client,
	}, nil
}
