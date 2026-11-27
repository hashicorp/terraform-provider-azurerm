package edgeactionversions

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EdgeActionVersionsClient struct {
	Client *resourcemanager.Client
}

func NewEdgeActionVersionsClientWithBaseURI(sdkApi sdkEnv.Api) (*EdgeActionVersionsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "edgeactionversions", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating EdgeActionVersionsClient: %+v", err)
	}

	return &EdgeActionVersionsClient{
		Client: client,
	}, nil
}
