package customizationtasks

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomizationTasksClient struct {
	Client *resourcemanager.Client
}

func NewCustomizationTasksClientWithBaseURI(sdkApi sdkEnv.Api) (*CustomizationTasksClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "customizationtasks", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating CustomizationTasksClient: %+v", err)
	}

	return &CustomizationTasksClient{
		Client: client,
	}, nil
}
