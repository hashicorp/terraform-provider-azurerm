package routingrules

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoutingRulesClient struct {
	Client *resourcemanager.Client
}

func NewRoutingRulesClientWithBaseURI(sdkApi sdkEnv.Api) (*RoutingRulesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "routingrules", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating RoutingRulesClient: %+v", err)
	}

	return &RoutingRulesClient{
		Client: client,
	}, nil
}
