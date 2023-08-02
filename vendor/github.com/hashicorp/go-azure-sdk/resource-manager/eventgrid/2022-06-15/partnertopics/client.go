package partnertopics

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PartnerTopicsClient struct {
	Client *resourcemanager.Client
}

func NewPartnerTopicsClientWithBaseURI(sdkApi sdkEnv.Api) (*PartnerTopicsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "partnertopics", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PartnerTopicsClient: %+v", err)
	}

	return &PartnerTopicsClient{
		Client: client,
	}, nil
}
