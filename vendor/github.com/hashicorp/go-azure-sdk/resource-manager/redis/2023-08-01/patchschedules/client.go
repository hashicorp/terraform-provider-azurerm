package patchschedules

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PatchSchedulesClient struct {
	Client *resourcemanager.Client
}

func NewPatchSchedulesClientWithBaseURI(sdkApi sdkEnv.Api) (*PatchSchedulesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "patchschedules", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PatchSchedulesClient: %+v", err)
	}

	return &PatchSchedulesClient{
		Client: client,
	}, nil
}
