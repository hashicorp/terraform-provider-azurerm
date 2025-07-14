package projectcatalogs

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProjectCatalogsClient struct {
	Client *resourcemanager.Client
}

func NewProjectCatalogsClientWithBaseURI(sdkApi sdkEnv.Api) (*ProjectCatalogsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "projectcatalogs", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ProjectCatalogsClient: %+v", err)
	}

	return &ProjectCatalogsClient{
		Client: client,
	}, nil
}
