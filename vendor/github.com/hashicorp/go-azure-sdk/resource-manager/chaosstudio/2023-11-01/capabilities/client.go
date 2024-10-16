package capabilities

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CapabilitiesClient struct {
	Client *resourcemanager.Client
}

func NewCapabilitiesClientWithBaseURI(sdkApi sdkEnv.Api) (*CapabilitiesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "capabilities", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating CapabilitiesClient: %+v", err)
	}

	return &CapabilitiesClient{
		Client: client,
	}, nil
}
