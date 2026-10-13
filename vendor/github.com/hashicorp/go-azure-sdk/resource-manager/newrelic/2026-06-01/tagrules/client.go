package tagrules

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TagRulesClient struct {
	Client *resourcemanager.Client
}

func NewTagRulesClientWithBaseURI(sdkApi sdkEnv.Api) (*TagRulesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "tagrules", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating TagRulesClient: %+v", err)
	}

	return &TagRulesClient{
		Client: client,
	}, nil
}
