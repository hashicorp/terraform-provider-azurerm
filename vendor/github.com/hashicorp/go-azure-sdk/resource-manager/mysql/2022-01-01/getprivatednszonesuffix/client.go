package getprivatednszonesuffix

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetPrivateDnsZoneSuffixClient struct {
	Client *resourcemanager.Client
}

func NewGetPrivateDnsZoneSuffixClientWithBaseURI(sdkApi sdkEnv.Api) (*GetPrivateDnsZoneSuffixClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "getprivatednszonesuffix", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating GetPrivateDnsZoneSuffixClient: %+v", err)
	}

	return &GetPrivateDnsZoneSuffixClient{
		Client: client,
	}, nil
}
