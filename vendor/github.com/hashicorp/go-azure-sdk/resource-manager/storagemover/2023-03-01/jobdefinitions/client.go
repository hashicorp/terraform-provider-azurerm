package jobdefinitions

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobDefinitionsClient struct {
	Client *resourcemanager.Client
}

func NewJobDefinitionsClientWithBaseURI(sdkApi sdkEnv.Api) (*JobDefinitionsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "jobdefinitions", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating JobDefinitionsClient: %+v", err)
	}

	return &JobDefinitionsClient{
		Client: client,
	}, nil
}
