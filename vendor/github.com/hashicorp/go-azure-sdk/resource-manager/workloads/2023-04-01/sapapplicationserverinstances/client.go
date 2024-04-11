package sapapplicationserverinstances

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SAPApplicationServerInstancesClient struct {
	Client *resourcemanager.Client
}

func NewSAPApplicationServerInstancesClientWithBaseURI(sdkApi sdkEnv.Api) (*SAPApplicationServerInstancesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "sapapplicationserverinstances", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SAPApplicationServerInstancesClient: %+v", err)
	}

	return &SAPApplicationServerInstancesClient{
		Client: client,
	}, nil
}
