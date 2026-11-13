package supercomputers

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SupercomputersClient struct {
	Client *resourcemanager.Client
}

func NewSupercomputersClientWithBaseURI(sdkApi sdkEnv.Api) (*SupercomputersClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "supercomputers", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SupercomputersClient: %+v", err)
	}

	return &SupercomputersClient{
		Client: client,
	}, nil
}
