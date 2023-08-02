package publishedblueprint

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PublishedBlueprintClient struct {
	Client *resourcemanager.Client
}

func NewPublishedBlueprintClientWithBaseURI(sdkApi sdkEnv.Api) (*PublishedBlueprintClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "publishedblueprint", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PublishedBlueprintClient: %+v", err)
	}

	return &PublishedBlueprintClient{
		Client: client,
	}, nil
}
