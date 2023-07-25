package snapshotpolicy

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SnapshotPolicyClient struct {
	Client *resourcemanager.Client
}

func NewSnapshotPolicyClientWithBaseURI(api environments.Api) (*SnapshotPolicyClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "snapshotpolicy", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SnapshotPolicyClient: %+v", err)
	}

	return &SnapshotPolicyClient{
		Client: client,
	}, nil
}
