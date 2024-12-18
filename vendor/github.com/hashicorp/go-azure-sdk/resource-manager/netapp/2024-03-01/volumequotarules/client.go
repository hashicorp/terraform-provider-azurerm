package volumequotarules

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VolumeQuotaRulesClient struct {
	Client *resourcemanager.Client
}

func NewVolumeQuotaRulesClientWithBaseURI(sdkApi sdkEnv.Api) (*VolumeQuotaRulesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "volumequotarules", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating VolumeQuotaRulesClient: %+v", err)
	}

	return &VolumeQuotaRulesClient{
		Client: client,
	}, nil
}
