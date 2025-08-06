package rolemanagementpolicies

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoleManagementPoliciesClient struct {
	Client *resourcemanager.Client
}

func NewRoleManagementPoliciesClientWithBaseURI(sdkApi sdkEnv.Api) (*RoleManagementPoliciesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "rolemanagementpolicies", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating RoleManagementPoliciesClient: %+v", err)
	}

	return &RoleManagementPoliciesClient{
		Client: client,
	}, nil
}
