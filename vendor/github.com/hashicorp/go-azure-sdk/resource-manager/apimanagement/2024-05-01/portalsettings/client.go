package portalsettings

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PortalSettingsClient struct {
	Client *resourcemanager.Client
}

func NewPortalSettingsClientWithBaseURI(sdkApi sdkEnv.Api) (*PortalSettingsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "portalsettings", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PortalSettingsClient: %+v", err)
	}

	return &PortalSettingsClient{
		Client: client,
	}, nil
}
