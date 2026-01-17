package scclusterrecords

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SCClusterRecordsClient struct {
	Client *resourcemanager.Client
}

func NewSCClusterRecordsClientWithBaseURI(sdkApi sdkEnv.Api) (*SCClusterRecordsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "scclusterrecords", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SCClusterRecordsClient: %+v", err)
	}

	return &SCClusterRecordsClient{
		Client: client,
	}, nil
}
