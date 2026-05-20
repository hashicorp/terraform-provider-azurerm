package settings

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SettingsClient struct {
	Client *resourcemanager.Client
}

func NewSettingsClientWithBaseURI(sdkApi sdkEnv.Api) (*SettingsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "settings", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SettingsClient: %+v", err)
	}

	return &SettingsClient{
		Client: client,
	}, nil
}
