package managementgroups

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagementGroupsClient struct {
	Client *resourcemanager.Client
}

func NewManagementGroupsClientWithBaseURI(sdkApi sdkEnv.Api) (*ManagementGroupsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "managementgroups", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ManagementGroupsClient: %+v", err)
	}

	return &ManagementGroupsClient{
		Client: client,
	}, nil
}
