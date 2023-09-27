package cacherules

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CacheRulesClient struct {
	Client *resourcemanager.Client
}

func NewCacheRulesClientWithBaseURI(sdkApi sdkEnv.Api) (*CacheRulesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "cacherules", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating CacheRulesClient: %+v", err)
	}

	return &CacheRulesClient{
		Client: client,
	}, nil
}
