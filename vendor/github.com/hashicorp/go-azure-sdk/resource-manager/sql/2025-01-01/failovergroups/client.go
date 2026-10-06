package failovergroups

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FailoverGroupsClient struct {
	Client *resourcemanager.Client
}

func NewFailoverGroupsClientWithBaseURI(sdkApi sdkEnv.Api) (*FailoverGroupsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "failovergroups", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating FailoverGroupsClient: %+v", err)
	}

	return &FailoverGroupsClient{
		Client: client,
	}, nil
}
