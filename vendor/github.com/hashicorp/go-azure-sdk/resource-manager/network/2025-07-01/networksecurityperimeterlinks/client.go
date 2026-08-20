package networksecurityperimeterlinks

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkSecurityPerimeterLinksClient struct {
	Client *resourcemanager.Client
}

func NewNetworkSecurityPerimeterLinksClientWithBaseURI(sdkApi sdkEnv.Api) (*NetworkSecurityPerimeterLinksClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "networksecurityperimeterlinks", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating NetworkSecurityPerimeterLinksClient: %+v", err)
	}

	return &NetworkSecurityPerimeterLinksClient{
		Client: client,
	}, nil
}
