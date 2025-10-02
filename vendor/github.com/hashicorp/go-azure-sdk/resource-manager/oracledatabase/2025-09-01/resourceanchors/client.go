package resourceanchors

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceAnchorsClient struct {
	Client *resourcemanager.Client
}

func NewResourceAnchorsClientWithBaseURI(sdkApi sdkEnv.Api) (*ResourceAnchorsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "resourceanchors", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ResourceAnchorsClient: %+v", err)
	}

	return &ResourceAnchorsClient{
		Client: client,
	}, nil
}
