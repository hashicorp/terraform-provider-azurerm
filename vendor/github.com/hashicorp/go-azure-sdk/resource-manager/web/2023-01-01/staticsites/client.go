package staticsites

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StaticSitesClient struct {
	Client *resourcemanager.Client
}

func NewStaticSitesClientWithBaseURI(sdkApi sdkEnv.Api) (*StaticSitesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "staticsites", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating StaticSitesClient: %+v", err)
	}

	return &StaticSitesClient{
		Client: client,
	}, nil
}
