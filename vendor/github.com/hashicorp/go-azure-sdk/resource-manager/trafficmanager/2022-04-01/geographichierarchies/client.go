package geographichierarchies

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GeographicHierarchiesClient struct {
	Client *resourcemanager.Client
}

func NewGeographicHierarchiesClientWithBaseURI(sdkApi sdkEnv.Api) (*GeographicHierarchiesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "geographichierarchies", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating GeographicHierarchiesClient: %+v", err)
	}

	return &GeographicHierarchiesClient{
		Client: client,
	}, nil
}
