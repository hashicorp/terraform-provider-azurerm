package autoimportjob

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoImportJobClient struct {
	Client *resourcemanager.Client
}

func NewAutoImportJobClientWithBaseURI(sdkApi sdkEnv.Api) (*AutoImportJobClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "autoimportjob", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AutoImportJobClient: %+v", err)
	}

	return &AutoImportJobClient{
		Client: client,
	}, nil
}
