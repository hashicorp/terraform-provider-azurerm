package ehnamespaces

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EHNamespacesClient struct {
	Client *resourcemanager.Client
}

func NewEHNamespacesClientWithBaseURI(sdkApi sdkEnv.Api) (*EHNamespacesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "ehnamespaces", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating EHNamespacesClient: %+v", err)
	}

	return &EHNamespacesClient{
		Client: client,
	}, nil
}
