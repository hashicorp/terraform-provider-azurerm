package resource

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceClient struct {
	Client *resourcemanager.Client
}

func NewResourceClientWithBaseURI(sdkApi sdkEnv.Api) (*ResourceClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "resource", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ResourceClient: %+v", err)
	}

	return &ResourceClient{
		Client: client,
	}, nil
}
