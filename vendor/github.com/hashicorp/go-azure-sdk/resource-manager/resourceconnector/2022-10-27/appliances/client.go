package appliances

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AppliancesClient struct {
	Client *resourcemanager.Client
}

func NewAppliancesClientWithBaseURI(sdkApi sdkEnv.Api) (*AppliancesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "appliances", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AppliancesClient: %+v", err)
	}

	return &AppliancesClient{
		Client: client,
	}, nil
}
