package sapcentralserverinstances

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SAPCentralServerInstancesClient struct {
	Client *resourcemanager.Client
}

func NewSAPCentralServerInstancesClientWithBaseURI(sdkApi sdkEnv.Api) (*SAPCentralServerInstancesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "sapcentralserverinstances", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SAPCentralServerInstancesClient: %+v", err)
	}

	return &SAPCentralServerInstancesClient{
		Client: client,
	}, nil
}
