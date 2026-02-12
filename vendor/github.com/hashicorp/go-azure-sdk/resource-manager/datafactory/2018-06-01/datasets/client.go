package datasets

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatasetsClient struct {
	Client *resourcemanager.Client
}

func NewDatasetsClientWithBaseURI(sdkApi sdkEnv.Api) (*DatasetsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "datasets", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DatasetsClient: %+v", err)
	}

	return &DatasetsClient{
		Client: client,
	}, nil
}
