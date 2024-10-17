package actiongroupsapis

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActionGroupsAPIsClient struct {
	Client *resourcemanager.Client
}

func NewActionGroupsAPIsClientWithBaseURI(sdkApi sdkEnv.Api) (*ActionGroupsAPIsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "actiongroupsapis", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ActionGroupsAPIsClient: %+v", err)
	}

	return &ActionGroupsAPIsClient{
		Client: client,
	}, nil
}
