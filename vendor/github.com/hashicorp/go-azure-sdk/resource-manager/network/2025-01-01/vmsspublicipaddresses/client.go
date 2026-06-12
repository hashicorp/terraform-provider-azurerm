package vmsspublicipaddresses

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VMSSPublicIPAddressesClient struct {
	Client *resourcemanager.Client
}

func NewVMSSPublicIPAddressesClientWithBaseURI(sdkApi sdkEnv.Api) (*VMSSPublicIPAddressesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "vmsspublicipaddresses", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating VMSSPublicIPAddressesClient: %+v", err)
	}

	return &VMSSPublicIPAddressesClient{
		Client: client,
	}, nil
}
