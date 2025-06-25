package autoexportjobs

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoExportJobsClient struct {
	Client *resourcemanager.Client
}

func NewAutoExportJobsClientWithBaseURI(sdkApi sdkEnv.Api) (*AutoExportJobsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "autoexportjobs", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AutoExportJobsClient: %+v", err)
	}

	return &AutoExportJobsClient{
		Client: client,
	}, nil
}
