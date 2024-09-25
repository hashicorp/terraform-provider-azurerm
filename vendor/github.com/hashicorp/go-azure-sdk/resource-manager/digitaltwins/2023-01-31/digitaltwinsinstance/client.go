package digitaltwinsinstance

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DigitalTwinsInstanceClient struct {
	Client *resourcemanager.Client
}

func NewDigitalTwinsInstanceClientWithBaseURI(sdkApi sdkEnv.Api) (*DigitalTwinsInstanceClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "digitaltwinsinstance", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DigitalTwinsInstanceClient: %+v", err)
	}

	return &DigitalTwinsInstanceClient{
		Client: client,
	}, nil
}
