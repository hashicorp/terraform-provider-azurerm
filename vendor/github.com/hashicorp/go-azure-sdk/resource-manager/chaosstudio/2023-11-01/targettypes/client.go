package targettypes

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TargetTypesClient struct {
	Client *resourcemanager.Client
}

func NewTargetTypesClientWithBaseURI(sdkApi sdkEnv.Api) (*TargetTypesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "targettypes", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating TargetTypesClient: %+v", err)
	}

	return &TargetTypesClient{
		Client: client,
	}, nil
}
