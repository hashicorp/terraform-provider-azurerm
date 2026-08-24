package projectconnectionresource

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProjectConnectionResourceClient struct {
	Client *resourcemanager.Client
}

func NewProjectConnectionResourceClientWithBaseURI(sdkApi sdkEnv.Api) (*ProjectConnectionResourceClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "projectconnectionresource", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ProjectConnectionResourceClient: %+v", err)
	}

	return &ProjectConnectionResourceClient{
		Client: client,
	}, nil
}
