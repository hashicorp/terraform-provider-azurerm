package share

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ShareClient struct {
	Client *resourcemanager.Client
}

func NewShareClientWithBaseURI(sdkApi sdkEnv.Api) (*ShareClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "share", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ShareClient: %+v", err)
	}

	return &ShareClient{
		Client: client,
	}, nil
}
