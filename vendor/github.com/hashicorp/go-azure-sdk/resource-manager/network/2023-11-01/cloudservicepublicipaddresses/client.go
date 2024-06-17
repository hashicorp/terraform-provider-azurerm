package cloudservicepublicipaddresses

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CloudServicePublicIPAddressesClient struct {
	Client *resourcemanager.Client
}

func NewCloudServicePublicIPAddressesClientWithBaseURI(sdkApi sdkEnv.Api) (*CloudServicePublicIPAddressesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "cloudservicepublicipaddresses", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating CloudServicePublicIPAddressesClient: %+v", err)
	}

	return &CloudServicePublicIPAddressesClient{
		Client: client,
	}, nil
}
