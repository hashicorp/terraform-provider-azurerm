package metadata

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MetadataClient struct {
	Client *resourcemanager.Client
}

func NewMetadataClientWithBaseURI(sdkApi sdkEnv.Api) (*MetadataClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "metadata", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating MetadataClient: %+v", err)
	}

	return &MetadataClient{
		Client: client,
	}, nil
}
