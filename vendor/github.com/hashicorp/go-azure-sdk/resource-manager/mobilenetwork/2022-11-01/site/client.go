package site

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SiteClient struct {
	Client *resourcemanager.Client
}

func NewSiteClientWithBaseURI(sdkApi sdkEnv.Api) (*SiteClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "site", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SiteClient: %+v", err)
	}

	return &SiteClient{
		Client: client,
	}, nil
}
