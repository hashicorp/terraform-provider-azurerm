package metadataschemas

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MetadataSchemasClient struct {
	Client *resourcemanager.Client
}

func NewMetadataSchemasClientWithBaseURI(sdkApi sdkEnv.Api) (*MetadataSchemasClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "metadataschemas", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating MetadataSchemasClient: %+v", err)
	}

	return &MetadataSchemasClient{
		Client: client,
	}, nil
}
