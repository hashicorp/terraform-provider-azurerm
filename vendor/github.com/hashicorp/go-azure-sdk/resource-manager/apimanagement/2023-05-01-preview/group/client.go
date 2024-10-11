package group

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GroupClient struct {
	Client *resourcemanager.Client
}

func NewGroupClientWithBaseURI(sdkApi sdkEnv.Api) (*GroupClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "group", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating GroupClient: %+v", err)
	}

	return &GroupClient{
		Client: client,
	}, nil
}
