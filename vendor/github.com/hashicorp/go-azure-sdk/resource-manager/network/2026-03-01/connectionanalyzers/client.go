package connectionanalyzers

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionAnalyzersClient struct {
	Client *resourcemanager.Client
}

func NewConnectionAnalyzersClientWithBaseURI(sdkApi sdkEnv.Api) (*ConnectionAnalyzersClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "connectionanalyzers", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ConnectionAnalyzersClient: %+v", err)
	}

	return &ConnectionAnalyzersClient{
		Client: client,
	}, nil
}
