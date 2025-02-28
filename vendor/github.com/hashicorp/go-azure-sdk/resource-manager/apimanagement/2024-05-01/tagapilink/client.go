package tagapilink

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TagApiLinkClient struct {
	Client *resourcemanager.Client
}

func NewTagApiLinkClientWithBaseURI(sdkApi sdkEnv.Api) (*TagApiLinkClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "tagapilink", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating TagApiLinkClient: %+v", err)
	}

	return &TagApiLinkClient{
		Client: client,
	}, nil
}
