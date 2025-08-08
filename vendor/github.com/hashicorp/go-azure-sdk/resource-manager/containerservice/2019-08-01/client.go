package v2019_08_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2019-08-01/agentpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2019-08-01/containerservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2019-08-01/managedclusters"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	AgentPools        *agentpools.AgentPoolsClient
	ContainerServices *containerservices.ContainerServicesClient
	ManagedClusters   *managedclusters.ManagedClustersClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	agentPoolsClient, err := agentpools.NewAgentPoolsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building AgentPools client: %+v", err)
	}
	configureFunc(agentPoolsClient.Client)

	containerServicesClient, err := containerservices.NewContainerServicesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ContainerServices client: %+v", err)
	}
	configureFunc(containerServicesClient.Client)

	managedClustersClient, err := managedclusters.NewManagedClustersClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ManagedClusters client: %+v", err)
	}
	configureFunc(managedClustersClient.Client)

	return &Client{
		AgentPools:        agentPoolsClient,
		ContainerServices: containerServicesClient,
		ManagedClusters:   managedClustersClient,
	}, nil
}
