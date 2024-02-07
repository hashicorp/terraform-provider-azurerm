package contact

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContactClient struct {
	Client *resourcemanager.Client
}

func NewContactClientWithBaseURI(sdkApi sdkEnv.Api) (*ContactClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "contact", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ContactClient: %+v", err)
	}

	return &ContactClient{
		Client: client,
	}, nil
}
