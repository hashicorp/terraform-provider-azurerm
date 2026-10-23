package groupquotalimits

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GroupQuotaLimitsClient struct {
	Client *resourcemanager.Client
}

func NewGroupQuotaLimitsClientWithBaseURI(sdkApi sdkEnv.Api) (*GroupQuotaLimitsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "groupquotalimits", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating GroupQuotaLimitsClient: %+v", err)
	}

	return &GroupQuotaLimitsClient{
		Client: client,
	}, nil
}
