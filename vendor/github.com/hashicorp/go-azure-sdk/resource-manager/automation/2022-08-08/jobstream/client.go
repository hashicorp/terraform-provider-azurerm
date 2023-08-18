package jobstream

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobStreamClient struct {
	Client *resourcemanager.Client
}

func NewJobStreamClientWithBaseURI(sdkApi sdkEnv.Api) (*JobStreamClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "jobstream", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating JobStreamClient: %+v", err)
	}

	return &JobStreamClient{
		Client: client,
	}, nil
}
