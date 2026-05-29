package v2019_06_01_preview

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2019-06-01-preview/agentpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2019-06-01-preview/registries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2019-06-01-preview/runs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2019-06-01-preview/taskruns"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2019-06-01-preview/tasks"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	AgentPools *agentpools.AgentPoolsClient
	Registries *registries.RegistriesClient
	Runs       *runs.RunsClient
	TaskRuns   *taskruns.TaskRunsClient
	Tasks      *tasks.TasksClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	agentPoolsClient, err := agentpools.NewAgentPoolsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building AgentPools client: %+v", err)
	}
	configureFunc(agentPoolsClient.Client)

	registriesClient, err := registries.NewRegistriesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Registries client: %+v", err)
	}
	configureFunc(registriesClient.Client)

	runsClient, err := runs.NewRunsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Runs client: %+v", err)
	}
	configureFunc(runsClient.Client)

	taskRunsClient, err := taskruns.NewTaskRunsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building TaskRuns client: %+v", err)
	}
	configureFunc(taskRunsClient.Client)

	tasksClient, err := tasks.NewTasksClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Tasks client: %+v", err)
	}
	configureFunc(tasksClient.Client)

	return &Client{
		AgentPools: agentPoolsClient,
		Registries: registriesClient,
		Runs:       runsClient,
		TaskRuns:   taskRunsClient,
		Tasks:      tasksClient,
	}, nil
}
