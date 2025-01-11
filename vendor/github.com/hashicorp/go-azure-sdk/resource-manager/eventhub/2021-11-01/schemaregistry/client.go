package schemaregistry

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SchemaRegistryClient struct {
	Client *resourcemanager.Client
}

func NewSchemaRegistryClientWithBaseURI(sdkApi sdkEnv.Api) (*SchemaRegistryClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "schemaregistry", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SchemaRegistryClient: %+v", err)
	}

	return &SchemaRegistryClient{
		Client: client,
	}, nil
}
