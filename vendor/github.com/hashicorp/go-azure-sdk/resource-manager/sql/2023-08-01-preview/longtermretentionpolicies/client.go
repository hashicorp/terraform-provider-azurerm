package longtermretentionpolicies

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LongTermRetentionPoliciesClient struct {
	Client *resourcemanager.Client
}

func NewLongTermRetentionPoliciesClientWithBaseURI(sdkApi sdkEnv.Api) (*LongTermRetentionPoliciesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "longtermretentionpolicies", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating LongTermRetentionPoliciesClient: %+v", err)
	}

	return &LongTermRetentionPoliciesClient{
		Client: client,
	}, nil
}
