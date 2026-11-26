package bigdatapools

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BigDataPoolsClient struct {
	Client *resourcemanager.Client
}

func NewBigDataPoolsClientWithBaseURI(sdkApi sdkEnv.Api) (*BigDataPoolsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "bigdatapools", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating BigDataPoolsClient: %+v", err)
	}

	return &BigDataPoolsClient{
		Client: client,
	}, nil
}
