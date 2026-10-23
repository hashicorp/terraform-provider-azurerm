package subscriptionquotaallocation

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubscriptionQuotaAllocationClient struct {
	Client *resourcemanager.Client
}

func NewSubscriptionQuotaAllocationClientWithBaseURI(sdkApi sdkEnv.Api) (*SubscriptionQuotaAllocationClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "subscriptionquotaallocation", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SubscriptionQuotaAllocationClient: %+v", err)
	}

	return &SubscriptionQuotaAllocationClient{
		Client: client,
	}, nil
}
