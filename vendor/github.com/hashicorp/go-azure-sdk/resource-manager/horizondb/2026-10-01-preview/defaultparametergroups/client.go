package defaultparametergroups

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DefaultParameterGroupsClient struct {
	Client *resourcemanager.Client
}

func NewDefaultParameterGroupsClientWithBaseURI(sdkApi sdkEnv.Api) (*DefaultParameterGroupsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "defaultparametergroups", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DefaultParameterGroupsClient: %+v", err)
	}

	return &DefaultParameterGroupsClient{
		Client: client,
	}, nil
}
