package sapsupportedsku

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SAPSupportedSkuClient struct {
	Client *resourcemanager.Client
}

func NewSAPSupportedSkuClientWithBaseURI(sdkApi sdkEnv.Api) (*SAPSupportedSkuClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "sapsupportedsku", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SAPSupportedSkuClient: %+v", err)
	}

	return &SAPSupportedSkuClient{
		Client: client,
	}, nil
}
