package accountcapabilityhost

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountCapabilityHostClient struct {
	Client *resourcemanager.Client
}

func NewAccountCapabilityHostClientWithBaseURI(sdkApi sdkEnv.Api) (*AccountCapabilityHostClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "accountcapabilityhost", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AccountCapabilityHostClient: %+v", err)
	}

	return &AccountCapabilityHostClient{
		Client: client,
	}, nil
}
