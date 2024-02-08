package accountfilters

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountFiltersClient struct {
	Client *resourcemanager.Client
}

func NewAccountFiltersClientWithBaseURI(sdkApi sdkEnv.Api) (*AccountFiltersClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "accountfilters", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AccountFiltersClient: %+v", err)
	}

	return &AccountFiltersClient{
		Client: client,
	}, nil
}
