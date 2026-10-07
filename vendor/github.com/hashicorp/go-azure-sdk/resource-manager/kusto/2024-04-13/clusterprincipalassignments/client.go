package clusterprincipalassignments

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterPrincipalAssignmentsClient struct {
	Client *resourcemanager.Client
}

func NewClusterPrincipalAssignmentsClientWithBaseURI(sdkApi sdkEnv.Api) (*ClusterPrincipalAssignmentsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "clusterprincipalassignments", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ClusterPrincipalAssignmentsClient: %+v", err)
	}

	return &ClusterPrincipalAssignmentsClient{
		Client: client,
	}, nil
}
