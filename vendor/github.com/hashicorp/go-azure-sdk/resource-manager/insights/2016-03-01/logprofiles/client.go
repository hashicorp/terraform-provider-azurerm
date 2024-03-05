package logprofiles

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LogProfilesClient struct {
	Client *resourcemanager.Client
}

func NewLogProfilesClientWithBaseURI(sdkApi sdkEnv.Api) (*LogProfilesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "logprofiles", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating LogProfilesClient: %+v", err)
	}

	return &LogProfilesClient{
		Client: client,
	}, nil
}
