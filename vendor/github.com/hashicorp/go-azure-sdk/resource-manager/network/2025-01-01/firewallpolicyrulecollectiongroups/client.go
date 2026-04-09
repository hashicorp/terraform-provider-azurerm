package firewallpolicyrulecollectiongroups

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FirewallPolicyRuleCollectionGroupsClient struct {
	Client *resourcemanager.Client
}

func NewFirewallPolicyRuleCollectionGroupsClientWithBaseURI(sdkApi sdkEnv.Api) (*FirewallPolicyRuleCollectionGroupsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "firewallpolicyrulecollectiongroups", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating FirewallPolicyRuleCollectionGroupsClient: %+v", err)
	}

	return &FirewallPolicyRuleCollectionGroupsClient{
		Client: client,
	}, nil
}
