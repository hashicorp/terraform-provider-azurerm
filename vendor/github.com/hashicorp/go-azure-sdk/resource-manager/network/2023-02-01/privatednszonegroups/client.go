package privatednszonegroups

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateDnsZoneGroupsClient struct {
	Client *resourcemanager.Client
}

func NewPrivateDnsZoneGroupsClientWithBaseURI(api environments.Api) (*PrivateDnsZoneGroupsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "privatednszonegroups", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PrivateDnsZoneGroupsClient: %+v", err)
	}

	return &PrivateDnsZoneGroupsClient{
		Client: client,
	}, nil
}
