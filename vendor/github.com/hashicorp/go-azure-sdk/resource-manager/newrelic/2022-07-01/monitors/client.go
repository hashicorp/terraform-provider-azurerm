package monitors

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MonitorsClient struct {
	Client *resourcemanager.Client
}

func NewMonitorsClientWithBaseURI(sdkApi sdkEnv.Api) (*MonitorsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "monitors", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating MonitorsClient: %+v", err)
	}

	return &MonitorsClient{
		Client: client,
	}, nil
}
