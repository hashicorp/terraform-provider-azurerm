package systemversions

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SystemVersionsClient struct {
	Client *resourcemanager.Client
}

func NewSystemVersionsClientWithBaseURI(sdkApi sdkEnv.Api) (*SystemVersionsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "systemversions", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SystemVersionsClient: %+v", err)
	}

	return &SystemVersionsClient{
		Client: client,
	}, nil
}
