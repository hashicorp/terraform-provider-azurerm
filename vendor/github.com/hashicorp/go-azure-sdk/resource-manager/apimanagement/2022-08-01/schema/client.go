package schema

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SchemaClient struct {
	Client *resourcemanager.Client
}

func NewSchemaClientWithBaseURI(sdkApi sdkEnv.Api) (*SchemaClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "schema", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SchemaClient: %+v", err)
	}

	return &SchemaClient{
		Client: client,
	}, nil
}
