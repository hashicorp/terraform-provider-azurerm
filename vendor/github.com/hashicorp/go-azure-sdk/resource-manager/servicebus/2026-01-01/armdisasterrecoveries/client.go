package armdisasterrecoveries

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ArmDisasterRecoveriesClient struct {
	Client *resourcemanager.Client
}

func NewArmDisasterRecoveriesClientWithBaseURI(sdkApi sdkEnv.Api) (*ArmDisasterRecoveriesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "armdisasterrecoveries", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ArmDisasterRecoveriesClient: %+v", err)
	}

	return &ArmDisasterRecoveriesClient{
		Client: client,
	}, nil
}
