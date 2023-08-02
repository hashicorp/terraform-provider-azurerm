package prerules

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PreRulesClient struct {
	Client *resourcemanager.Client
}

func NewPreRulesClientWithBaseURI(sdkApi sdkEnv.Api) (*PreRulesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "prerules", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PreRulesClient: %+v", err)
	}

	return &PreRulesClient{
		Client: client,
	}, nil
}
