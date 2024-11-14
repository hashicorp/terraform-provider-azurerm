package raicontentfilters

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RaiContentFiltersClient struct {
	Client *resourcemanager.Client
}

func NewRaiContentFiltersClientWithBaseURI(sdkApi sdkEnv.Api) (*RaiContentFiltersClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "raicontentfilters", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating RaiContentFiltersClient: %+v", err)
	}

	return &RaiContentFiltersClient{
		Client: client,
	}, nil
}
