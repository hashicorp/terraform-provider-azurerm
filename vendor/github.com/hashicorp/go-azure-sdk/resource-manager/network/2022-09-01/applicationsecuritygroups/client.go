package applicationsecuritygroups

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationSecurityGroupsClient struct {
	Client *resourcemanager.Client
}

func NewApplicationSecurityGroupsClientWithBaseURI(api environments.Api) (*ApplicationSecurityGroupsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "applicationsecuritygroups", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ApplicationSecurityGroupsClient: %+v", err)
	}

	return &ApplicationSecurityGroupsClient{
		Client: client,
	}, nil
}
