package giminorversions

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GiMinorVersionsClient struct {
	Client *resourcemanager.Client
}

func NewGiMinorVersionsClientWithBaseURI(sdkApi sdkEnv.Api) (*GiMinorVersionsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "giminorversions", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating GiMinorVersionsClient: %+v", err)
	}

	return &GiMinorVersionsClient{
		Client: client,
	}, nil
}
