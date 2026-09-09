package jobstepexecutions

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobStepExecutionsClient struct {
	Client *resourcemanager.Client
}

func NewJobStepExecutionsClientWithBaseURI(sdkApi sdkEnv.Api) (*JobStepExecutionsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "jobstepexecutions", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating JobStepExecutionsClient: %+v", err)
	}

	return &JobStepExecutionsClient{
		Client: client,
	}, nil
}
