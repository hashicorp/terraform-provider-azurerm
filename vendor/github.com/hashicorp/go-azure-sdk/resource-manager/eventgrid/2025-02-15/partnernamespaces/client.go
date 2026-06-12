package partnernamespaces

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PartnerNamespacesClient struct {
	Client *resourcemanager.Client
}

func NewPartnerNamespacesClientWithBaseURI(sdkApi sdkEnv.Api) (*PartnerNamespacesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "partnernamespaces", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PartnerNamespacesClient: %+v", err)
	}

	return &PartnerNamespacesClient{
		Client: client,
	}, nil
}
