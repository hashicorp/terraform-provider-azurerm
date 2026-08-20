package interconnectgroups

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InterconnectGroupsClient struct {
	Client *resourcemanager.Client
}

func NewInterconnectGroupsClientWithBaseURI(sdkApi sdkEnv.Api) (*InterconnectGroupsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "interconnectgroups", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating InterconnectGroupsClient: %+v", err)
	}

	return &InterconnectGroupsClient{
		Client: client,
	}, nil
}
