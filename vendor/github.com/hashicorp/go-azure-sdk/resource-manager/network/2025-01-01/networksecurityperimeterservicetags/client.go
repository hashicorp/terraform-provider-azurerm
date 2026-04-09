package networksecurityperimeterservicetags

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkSecurityPerimeterServiceTagsClient struct {
	Client *resourcemanager.Client
}

func NewNetworkSecurityPerimeterServiceTagsClientWithBaseURI(sdkApi sdkEnv.Api) (*NetworkSecurityPerimeterServiceTagsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "networksecurityperimeterservicetags", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating NetworkSecurityPerimeterServiceTagsClient: %+v", err)
	}

	return &NetworkSecurityPerimeterServiceTagsClient{
		Client: client,
	}, nil
}
