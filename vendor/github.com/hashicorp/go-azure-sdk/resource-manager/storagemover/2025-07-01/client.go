package v2025_07_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/agents"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/endpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/jobdefinitions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/jobruns"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/projects"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/storagemovers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	Agents         *agents.AgentsClient
	Endpoints      *endpoints.EndpointsClient
	JobDefinitions *jobdefinitions.JobDefinitionsClient
	JobRuns        *jobruns.JobRunsClient
	Projects       *projects.ProjectsClient
	StorageMovers  *storagemovers.StorageMoversClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	agentsClient, err := agents.NewAgentsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Agents client: %+v", err)
	}
	configureFunc(agentsClient.Client)

	endpointsClient, err := endpoints.NewEndpointsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Endpoints client: %+v", err)
	}
	configureFunc(endpointsClient.Client)

	jobDefinitionsClient, err := jobdefinitions.NewJobDefinitionsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building JobDefinitions client: %+v", err)
	}
	configureFunc(jobDefinitionsClient.Client)

	jobRunsClient, err := jobruns.NewJobRunsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building JobRuns client: %+v", err)
	}
	configureFunc(jobRunsClient.Client)

	projectsClient, err := projects.NewProjectsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Projects client: %+v", err)
	}
	configureFunc(projectsClient.Client)

	storageMoversClient, err := storagemovers.NewStorageMoversClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building StorageMovers client: %+v", err)
	}
	configureFunc(storageMoversClient.Client)

	return &Client{
		Agents:         agentsClient,
		Endpoints:      endpointsClient,
		JobDefinitions: jobDefinitionsClient,
		JobRuns:        jobRunsClient,
		Projects:       projectsClient,
		StorageMovers:  storageMoversClient,
	}, nil
}
