package quotabycounterkeys

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QuotaByCounterKeysClient struct {
	Client *resourcemanager.Client
}

func NewQuotaByCounterKeysClientWithBaseURI(sdkApi sdkEnv.Api) (*QuotaByCounterKeysClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "quotabycounterkeys", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating QuotaByCounterKeysClient: %+v", err)
	}

	return &QuotaByCounterKeysClient{
		Client: client,
	}, nil
}
