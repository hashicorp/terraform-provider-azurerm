package networkgroups

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkGroupsClient struct {
	Client *resourcemanager.Client
}

func NewNetworkGroupsClientWithBaseURI(api environments.Api) (*NetworkGroupsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "networkgroups", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating NetworkGroupsClient: %+v", err)
	}

	return &NetworkGroupsClient{
		Client: client,
	}, nil
}
