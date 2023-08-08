package queuesauthorizationrule

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueuesAuthorizationRuleClient struct {
	Client *resourcemanager.Client
}

func NewQueuesAuthorizationRuleClientWithBaseURI(sdkApi sdkEnv.Api) (*QueuesAuthorizationRuleClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "queuesauthorizationrule", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating QueuesAuthorizationRuleClient: %+v", err)
	}

	return &QueuesAuthorizationRuleClient{
		Client: client,
	}, nil
}
