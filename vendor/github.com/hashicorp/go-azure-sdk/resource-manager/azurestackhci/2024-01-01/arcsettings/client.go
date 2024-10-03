package arcsettings

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ArcSettingsClient struct {
	Client *resourcemanager.Client
}

func NewArcSettingsClientWithBaseURI(sdkApi sdkEnv.Api) (*ArcSettingsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "arcsettings", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ArcSettingsClient: %+v", err)
	}

	return &ArcSettingsClient{
		Client: client,
	}, nil
}
