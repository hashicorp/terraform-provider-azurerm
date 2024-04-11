package profiles

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProfilesClient struct {
	Client *resourcemanager.Client
}

func NewProfilesClientWithBaseURI(sdkApi sdkEnv.Api) (*ProfilesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "profiles", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ProfilesClient: %+v", err)
	}

	return &ProfilesClient{
		Client: client,
	}, nil
}
