package creators

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CreatorsClient struct {
	Client *resourcemanager.Client
}

func NewCreatorsClientWithBaseURI(sdkApi sdkEnv.Api) (*CreatorsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "creators", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating CreatorsClient: %+v", err)
	}

	return &CreatorsClient{
		Client: client,
	}, nil
}
