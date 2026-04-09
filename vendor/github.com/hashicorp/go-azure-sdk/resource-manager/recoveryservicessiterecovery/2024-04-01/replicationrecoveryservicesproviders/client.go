package replicationrecoveryservicesproviders

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplicationRecoveryServicesProvidersClient struct {
	Client *resourcemanager.Client
}

func NewReplicationRecoveryServicesProvidersClientWithBaseURI(sdkApi sdkEnv.Api) (*ReplicationRecoveryServicesProvidersClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "replicationrecoveryservicesproviders", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ReplicationRecoveryServicesProvidersClient: %+v", err)
	}

	return &ReplicationRecoveryServicesProvidersClient{
		Client: client,
	}, nil
}
