package retentionpolicies

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RetentionPoliciesClient struct {
	Client *resourcemanager.Client
}

func NewRetentionPoliciesClientWithBaseURI(sdkApi sdkEnv.Api) (*RetentionPoliciesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "retentionpolicies", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating RetentionPoliciesClient: %+v", err)
	}

	return &RetentionPoliciesClient{
		Client: client,
	}, nil
}
