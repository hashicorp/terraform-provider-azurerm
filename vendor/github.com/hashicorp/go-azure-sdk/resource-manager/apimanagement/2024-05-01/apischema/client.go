package apischema

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiSchemaClient struct {
	Client *resourcemanager.Client
}

func NewApiSchemaClientWithBaseURI(sdkApi sdkEnv.Api) (*ApiSchemaClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "apischema", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ApiSchemaClient: %+v", err)
	}

	return &ApiSchemaClient{
		Client: client,
	}, nil
}
