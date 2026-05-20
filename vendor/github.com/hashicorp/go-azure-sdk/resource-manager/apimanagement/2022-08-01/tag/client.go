package tag

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TagClient struct {
	Client *resourcemanager.Client
}

func NewTagClientWithBaseURI(sdkApi sdkEnv.Api) (*TagClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "tag", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating TagClient: %+v", err)
	}

	return &TagClient{
		Client: client,
	}, nil
}
