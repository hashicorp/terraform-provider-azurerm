package configurationassignments

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigurationAssignmentsClient struct {
	Client *resourcemanager.Client
}

func NewConfigurationAssignmentsClientWithBaseURI(sdkApi sdkEnv.Api) (*ConfigurationAssignmentsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "configurationassignments", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ConfigurationAssignmentsClient: %+v", err)
	}

	return &ConfigurationAssignmentsClient{
		Client: client,
	}, nil
}
