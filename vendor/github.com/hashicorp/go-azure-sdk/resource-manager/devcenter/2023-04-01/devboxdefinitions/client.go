package devboxdefinitions

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DevBoxDefinitionsClient struct {
	Client *resourcemanager.Client
}

func NewDevBoxDefinitionsClientWithBaseURI(sdkApi sdkEnv.Api) (*DevBoxDefinitionsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "devboxdefinitions", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DevBoxDefinitionsClient: %+v", err)
	}

	return &DevBoxDefinitionsClient{
		Client: client,
	}, nil
}
