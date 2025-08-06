package experiments

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExperimentsClient struct {
	Client *resourcemanager.Client
}

func NewExperimentsClientWithBaseURI(sdkApi sdkEnv.Api) (*ExperimentsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "experiments", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ExperimentsClient: %+v", err)
	}

	return &ExperimentsClient{
		Client: client,
	}, nil
}
