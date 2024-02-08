package workbooksapis

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkbooksAPIsClient struct {
	Client *resourcemanager.Client
}

func NewWorkbooksAPIsClientWithBaseURI(sdkApi sdkEnv.Api) (*WorkbooksAPIsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "workbooksapis", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating WorkbooksAPIsClient: %+v", err)
	}

	return &WorkbooksAPIsClient{
		Client: client,
	}, nil
}
