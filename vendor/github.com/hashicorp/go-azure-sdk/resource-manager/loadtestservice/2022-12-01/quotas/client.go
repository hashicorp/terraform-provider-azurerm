package quotas

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QuotasClient struct {
	Client *resourcemanager.Client
}

func NewQuotasClientWithBaseURI(sdkApi sdkEnv.Api) (*QuotasClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "quotas", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating QuotasClient: %+v", err)
	}

	return &QuotasClient{
		Client: client,
	}, nil
}
