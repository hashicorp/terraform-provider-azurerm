package managementpolicies

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagementPoliciesClient struct {
	Client *resourcemanager.Client
}

func NewManagementPoliciesClientWithBaseURI(sdkApi sdkEnv.Api) (*ManagementPoliciesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "managementpolicies", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ManagementPoliciesClient: %+v", err)
	}

	return &ManagementPoliciesClient{
		Client: client,
	}, nil
}
