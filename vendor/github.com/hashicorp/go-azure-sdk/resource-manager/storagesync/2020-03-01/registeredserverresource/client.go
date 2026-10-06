package registeredserverresource

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegisteredServerResourceClient struct {
	Client *resourcemanager.Client
}

func NewRegisteredServerResourceClientWithBaseURI(sdkApi sdkEnv.Api) (*RegisteredServerResourceClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "registeredserverresource", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating RegisteredServerResourceClient: %+v", err)
	}

	return &RegisteredServerResourceClient{
		Client: client,
	}, nil
}
