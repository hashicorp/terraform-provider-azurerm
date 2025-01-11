package archiveversions

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ArchiveVersionsClient struct {
	Client *resourcemanager.Client
}

func NewArchiveVersionsClientWithBaseURI(sdkApi sdkEnv.Api) (*ArchiveVersionsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "archiveversions", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ArchiveVersionsClient: %+v", err)
	}

	return &ArchiveVersionsClient{
		Client: client,
	}, nil
}
