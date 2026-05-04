package networksecurityperimeteraccessrules

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkSecurityPerimeterAccessRulesClient struct {
	Client *resourcemanager.Client
}

func NewNetworkSecurityPerimeterAccessRulesClientWithBaseURI(sdkApi sdkEnv.Api) (*NetworkSecurityPerimeterAccessRulesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "networksecurityperimeteraccessrules", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating NetworkSecurityPerimeterAccessRulesClient: %+v", err)
	}

	return &NetworkSecurityPerimeterAccessRulesClient{
		Client: client,
	}, nil
}
