package projects

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProjectsClient struct {
	Client *resourcemanager.Client
}

func NewProjectsClientWithBaseURI(sdkApi sdkEnv.Api) (*ProjectsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "projects", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ProjectsClient: %+v", err)
	}

	return &ProjectsClient{
		Client: client,
	}, nil
}
