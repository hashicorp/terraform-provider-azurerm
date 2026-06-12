package daprcomponents

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DaprComponentsClient struct {
	Client *resourcemanager.Client
}

func NewDaprComponentsClientWithBaseURI(sdkApi sdkEnv.Api) (*DaprComponentsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "daprcomponents", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DaprComponentsClient: %+v", err)
	}

	return &DaprComponentsClient{
		Client: client,
	}, nil
}
