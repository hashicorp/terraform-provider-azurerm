package scopemaps

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScopeMapsClient struct {
	Client *resourcemanager.Client
}

func NewScopeMapsClientWithBaseURI(sdkApi sdkEnv.Api) (*ScopeMapsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "scopemaps", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ScopeMapsClient: %+v", err)
	}

	return &ScopeMapsClient{
		Client: client,
	}, nil
}
