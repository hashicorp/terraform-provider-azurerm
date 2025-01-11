package objectreplicationpolicies

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ObjectReplicationPoliciesClient struct {
	Client *resourcemanager.Client
}

func NewObjectReplicationPoliciesClientWithBaseURI(sdkApi sdkEnv.Api) (*ObjectReplicationPoliciesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "objectreplicationpolicies", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ObjectReplicationPoliciesClient: %+v", err)
	}

	return &ObjectReplicationPoliciesClient{
		Client: client,
	}, nil
}
