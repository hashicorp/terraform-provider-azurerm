package partnerregistrations

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PartnerRegistrationsClient struct {
	Client *resourcemanager.Client
}

func NewPartnerRegistrationsClientWithBaseURI(sdkApi sdkEnv.Api) (*PartnerRegistrationsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "partnerregistrations", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PartnerRegistrationsClient: %+v", err)
	}

	return &PartnerRegistrationsClient{
		Client: client,
	}, nil
}
