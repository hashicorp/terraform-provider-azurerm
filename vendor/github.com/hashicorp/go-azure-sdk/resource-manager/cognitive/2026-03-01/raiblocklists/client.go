package raiblocklists

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RaiBlocklistsClient struct {
	Client *resourcemanager.Client
}

func NewRaiBlocklistsClientWithBaseURI(sdkApi sdkEnv.Api) (*RaiBlocklistsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "raiblocklists", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating RaiBlocklistsClient: %+v", err)
	}

	return &RaiBlocklistsClient{
		Client: client,
	}, nil
}
