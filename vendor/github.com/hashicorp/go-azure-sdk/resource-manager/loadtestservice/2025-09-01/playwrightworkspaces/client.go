package playwrightworkspaces

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PlaywrightWorkspacesClient struct {
	Client *resourcemanager.Client
}

func NewPlaywrightWorkspacesClientWithBaseURI(sdkApi sdkEnv.Api) (*PlaywrightWorkspacesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "playwrightworkspaces", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PlaywrightWorkspacesClient: %+v", err)
	}

	return &PlaywrightWorkspacesClient{
		Client: client,
	}, nil
}
