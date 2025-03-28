package galleries

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GalleriesClient struct {
	Client *resourcemanager.Client
}

func NewGalleriesClientWithBaseURI(sdkApi sdkEnv.Api) (*GalleriesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "galleries", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating GalleriesClient: %+v", err)
	}

	return &GalleriesClient{
		Client: client,
	}, nil
}
