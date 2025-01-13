package virtualappliancesites

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualApplianceSitesClient struct {
	Client *resourcemanager.Client
}

func NewVirtualApplianceSitesClientWithBaseURI(sdkApi sdkEnv.Api) (*VirtualApplianceSitesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "virtualappliancesites", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating VirtualApplianceSitesClient: %+v", err)
	}

	return &VirtualApplianceSitesClient{
		Client: client,
	}, nil
}
