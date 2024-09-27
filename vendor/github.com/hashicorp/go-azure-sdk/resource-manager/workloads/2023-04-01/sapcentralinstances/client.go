package sapcentralinstances

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SAPCentralInstancesClient struct {
	Client *resourcemanager.Client
}

func NewSAPCentralInstancesClientWithBaseURI(sdkApi sdkEnv.Api) (*SAPCentralInstancesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "sapcentralinstances", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SAPCentralInstancesClient: %+v", err)
	}

	return &SAPCentralInstancesClient{
		Client: client,
	}, nil
}
