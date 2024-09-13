package replicationpolicies

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplicationPoliciesClient struct {
	Client *resourcemanager.Client
}

func NewReplicationPoliciesClientWithBaseURI(sdkApi sdkEnv.Api) (*ReplicationPoliciesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "replicationpolicies", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ReplicationPoliciesClient: %+v", err)
	}

	return &ReplicationPoliciesClient{
		Client: client,
	}, nil
}
