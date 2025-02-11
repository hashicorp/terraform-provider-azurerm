package routingrulecollections

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoutingRuleCollectionsClient struct {
	Client *resourcemanager.Client
}

func NewRoutingRuleCollectionsClientWithBaseURI(sdkApi sdkEnv.Api) (*RoutingRuleCollectionsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "routingrulecollections", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating RoutingRuleCollectionsClient: %+v", err)
	}

	return &RoutingRuleCollectionsClient{
		Client: client,
	}, nil
}
