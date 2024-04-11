package streamingjobs

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StreamingJobsClient struct {
	Client *resourcemanager.Client
}

func NewStreamingJobsClientWithBaseURI(sdkApi sdkEnv.Api) (*StreamingJobsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "streamingjobs", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating StreamingJobsClient: %+v", err)
	}

	return &StreamingJobsClient{
		Client: client,
	}, nil
}
