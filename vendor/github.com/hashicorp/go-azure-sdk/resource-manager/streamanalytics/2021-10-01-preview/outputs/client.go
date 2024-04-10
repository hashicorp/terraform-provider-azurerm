package outputs

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OutputsClient struct {
	Client *resourcemanager.Client
}

func NewOutputsClientWithBaseURI(sdkApi sdkEnv.Api) (*OutputsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "outputs", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating OutputsClient: %+v", err)
	}

	return &OutputsClient{
		Client: client,
	}, nil
}
