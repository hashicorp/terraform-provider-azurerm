package workbooktemplatesapis

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkbookTemplatesAPIsClient struct {
	Client *resourcemanager.Client
}

func NewWorkbookTemplatesAPIsClientWithBaseURI(api environments.Api) (*WorkbookTemplatesAPIsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "workbooktemplatesapis", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating WorkbookTemplatesAPIsClient: %+v", err)
	}

	return &WorkbookTemplatesAPIsClient{
		Client: client,
	}, nil
}
