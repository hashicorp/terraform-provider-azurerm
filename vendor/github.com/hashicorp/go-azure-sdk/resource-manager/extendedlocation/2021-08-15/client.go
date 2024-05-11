package v2021_08_15

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	CustomLocations *customlocations.CustomLocationsClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	customLocationsClient, err := customlocations.NewCustomLocationsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building CustomLocations client: %+v", err)
	}
	configureFunc(customLocationsClient.Client)

	return &Client{
		CustomLocations: customLocationsClient,
	}, nil
}
