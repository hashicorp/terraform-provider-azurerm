package geocatalogs

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GeoCatalogsClient struct {
	Client *resourcemanager.Client
}

func NewGeoCatalogsClientWithBaseURI(sdkApi sdkEnv.Api) (*GeoCatalogsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "geocatalogs", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating GeoCatalogsClient: %+v", err)
	}

	return &GeoCatalogsClient{
		Client: client,
	}, nil
}
