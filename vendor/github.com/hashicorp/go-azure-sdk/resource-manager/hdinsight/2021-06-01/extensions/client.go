package extensions

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExtensionsClient struct {
	Client *resourcemanager.Client
}

func NewExtensionsClientWithBaseURI(sdkApi sdkEnv.Api) (*ExtensionsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "extensions", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ExtensionsClient: %+v", err)
	}

	return &ExtensionsClient{
		Client: client,
	}, nil
}
