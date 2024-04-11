package sapvirtualinstances

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SAPVirtualInstancesClient struct {
	Client *resourcemanager.Client
}

func NewSAPVirtualInstancesClientWithBaseURI(sdkApi sdkEnv.Api) (*SAPVirtualInstancesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "sapvirtualinstances", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SAPVirtualInstancesClient: %+v", err)
	}

	return &SAPVirtualInstancesClient{
		Client: client,
	}, nil
}
