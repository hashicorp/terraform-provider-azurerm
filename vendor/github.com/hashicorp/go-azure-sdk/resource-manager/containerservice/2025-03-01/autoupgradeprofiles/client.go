package autoupgradeprofiles

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoUpgradeProfilesClient struct {
	Client *resourcemanager.Client
}

func NewAutoUpgradeProfilesClientWithBaseURI(sdkApi sdkEnv.Api) (*AutoUpgradeProfilesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "autoupgradeprofiles", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AutoUpgradeProfilesClient: %+v", err)
	}

	return &AutoUpgradeProfilesClient{
		Client: client,
	}, nil
}
