package topicsauthorizationrule

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TopicsAuthorizationRuleClient struct {
	Client *resourcemanager.Client
}

func NewTopicsAuthorizationRuleClientWithBaseURI(sdkApi sdkEnv.Api) (*TopicsAuthorizationRuleClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "topicsauthorizationrule", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating TopicsAuthorizationRuleClient: %+v", err)
	}

	return &TopicsAuthorizationRuleClient{
		Client: client,
	}, nil
}
