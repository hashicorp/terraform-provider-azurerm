package loadbalancers

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LoadBalancersClient struct {
	Client *resourcemanager.Client
}

func NewLoadBalancersClientWithBaseURI(sdkApi sdkEnv.Api) (*LoadBalancersClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "loadbalancers", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating LoadBalancersClient: %+v", err)
	}

	return &LoadBalancersClient{
		Client: client,
	}, nil
}
