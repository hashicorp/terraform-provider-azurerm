package schemagroups

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SchemaGroupsClient struct {
	Client *resourcemanager.Client
}

func NewSchemaGroupsClientWithBaseURI(sdkApi sdkEnv.Api) (*SchemaGroupsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "schemagroups", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SchemaGroupsClient: %+v", err)
	}

	return &SchemaGroupsClient{
		Client: client,
	}, nil
}
