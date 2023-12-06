package ddosprotectionplans

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DdosProtectionPlansClient struct {
	Client *resourcemanager.Client
}

func NewDdosProtectionPlansClientWithBaseURI(sdkApi sdkEnv.Api) (*DdosProtectionPlansClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "ddosprotectionplans", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DdosProtectionPlansClient: %+v", err)
	}

	return &DdosProtectionPlansClient{
		Client: client,
	}, nil
}
