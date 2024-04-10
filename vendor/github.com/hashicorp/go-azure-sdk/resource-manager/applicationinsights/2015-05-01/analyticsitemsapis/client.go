package analyticsitemsapis

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AnalyticsItemsAPIsClient struct {
	Client *resourcemanager.Client
}

func NewAnalyticsItemsAPIsClientWithBaseURI(sdkApi sdkEnv.Api) (*AnalyticsItemsAPIsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "analyticsitemsapis", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AnalyticsItemsAPIsClient: %+v", err)
	}

	return &AnalyticsItemsAPIsClient{
		Client: client,
	}, nil
}
