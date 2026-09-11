package schedulers

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SchedulersClient struct {
	Client *resourcemanager.Client
}

func NewSchedulersClientWithBaseURI(sdkApi sdkEnv.Api) (*SchedulersClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "schedulers", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SchedulersClient: %+v", err)
	}

	return &SchedulersClient{
		Client: client,
	}, nil
}
