package serverbasedperformancetier

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerBasedPerformanceTierClient struct {
	Client *resourcemanager.Client
}

func NewServerBasedPerformanceTierClientWithBaseURI(sdkApi sdkEnv.Api) (*ServerBasedPerformanceTierClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "serverbasedperformancetier", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ServerBasedPerformanceTierClient: %+v", err)
	}

	return &ServerBasedPerformanceTierClient{
		Client: client,
	}, nil
}
