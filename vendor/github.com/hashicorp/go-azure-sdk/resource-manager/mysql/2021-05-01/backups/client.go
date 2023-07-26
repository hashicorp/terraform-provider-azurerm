package backups

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupsClient struct {
	Client *resourcemanager.Client
}

func NewBackupsClientWithBaseURI(api environments.Api) (*BackupsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "backups", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating BackupsClient: %+v", err)
	}

	return &BackupsClient{
		Client: client,
	}, nil
}
