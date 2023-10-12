package applications

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationsClient struct {
	Client *resourcemanager.Client
}

func NewApplicationsClientWithBaseURI(sdkApi sdkEnv.Api) (*ApplicationsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "applications", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ApplicationsClient: %+v", err)
	}

	return &ApplicationsClient{
		Client: client,
	}, nil
}
