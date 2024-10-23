package zones

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ZonesClient struct {
	Client *resourcemanager.Client
}

func NewZonesClientWithBaseURI(sdkApi sdkEnv.Api) (*ZonesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "zones", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ZonesClient: %+v", err)
	}

	return &ZonesClient{
		Client: client,
	}, nil
}
