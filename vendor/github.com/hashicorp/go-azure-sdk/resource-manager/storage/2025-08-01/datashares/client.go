package datashares

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataSharesClient struct {
	Client *resourcemanager.Client
}

func NewDataSharesClientWithBaseURI(sdkApi sdkEnv.Api) (*DataSharesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "datashares", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DataSharesClient: %+v", err)
	}

	return &DataSharesClient{
		Client: client,
	}, nil
}
