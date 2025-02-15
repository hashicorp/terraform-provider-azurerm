package appserviceplans

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AppServicePlansClient struct {
	Client *resourcemanager.Client
}

func NewAppServicePlansClientWithBaseURI(sdkApi sdkEnv.Api) (*AppServicePlansClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "appserviceplans", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AppServicePlansClient: %+v", err)
	}

	return &AppServicePlansClient{
		Client: client,
	}, nil
}
