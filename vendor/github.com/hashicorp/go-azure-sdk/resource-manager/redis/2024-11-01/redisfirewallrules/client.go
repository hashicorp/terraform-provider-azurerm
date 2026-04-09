package redisfirewallrules

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RedisFirewallRulesClient struct {
	Client *resourcemanager.Client
}

func NewRedisFirewallRulesClientWithBaseURI(sdkApi sdkEnv.Api) (*RedisFirewallRulesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "redisfirewallrules", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating RedisFirewallRulesClient: %+v", err)
	}

	return &RedisFirewallRulesClient{
		Client: client,
	}, nil
}
