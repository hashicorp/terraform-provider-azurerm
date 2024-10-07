package usagemodels

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UsageModelsClient struct {
	Client *resourcemanager.Client
}

func NewUsageModelsClientWithBaseURI(sdkApi sdkEnv.Api) (*UsageModelsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "usagemodels", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating UsageModelsClient: %+v", err)
	}

	return &UsageModelsClient{
		Client: client,
	}, nil
}
