package scopeconnections

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScopeConnectionsClient struct {
	Client *resourcemanager.Client
}

func NewScopeConnectionsClientWithBaseURI(sdkApi sdkEnv.Api) (*ScopeConnectionsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "scopeconnections", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ScopeConnectionsClient: %+v", err)
	}

	return &ScopeConnectionsClient{
		Client: client,
	}, nil
}
