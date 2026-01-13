package backupsautomaticandondemand

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupsAutomaticAndOnDemandClient struct {
	Client *resourcemanager.Client
}

func NewBackupsAutomaticAndOnDemandClientWithBaseURI(sdkApi sdkEnv.Api) (*BackupsAutomaticAndOnDemandClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "backupsautomaticandondemand", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating BackupsAutomaticAndOnDemandClient: %+v", err)
	}

	return &BackupsAutomaticAndOnDemandClient{
		Client: client,
	}, nil
}
