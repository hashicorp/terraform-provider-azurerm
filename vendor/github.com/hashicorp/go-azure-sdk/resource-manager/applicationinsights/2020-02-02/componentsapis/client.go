package componentsapis

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ComponentsAPIsClient struct {
	Client *resourcemanager.Client
}

func NewComponentsAPIsClientWithBaseURI(sdkApi sdkEnv.Api) (*ComponentsAPIsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "componentsapis", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ComponentsAPIsClient: %+v", err)
	}

	return &ComponentsAPIsClient{
		Client: client,
	}, nil
}
