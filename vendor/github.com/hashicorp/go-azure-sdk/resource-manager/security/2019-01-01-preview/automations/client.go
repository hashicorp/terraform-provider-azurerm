package automations

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutomationsClient struct {
	Client *resourcemanager.Client
}

func NewAutomationsClientWithBaseURI(sdkApi sdkEnv.Api) (*AutomationsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "automations", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AutomationsClient: %+v", err)
	}

	return &AutomationsClient{
		Client: client,
	}, nil
}
