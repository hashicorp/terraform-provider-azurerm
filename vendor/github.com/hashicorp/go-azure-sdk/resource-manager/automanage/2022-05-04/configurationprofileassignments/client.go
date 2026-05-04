package configurationprofileassignments

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigurationProfileAssignmentsClient struct {
	Client *resourcemanager.Client
}

func NewConfigurationProfileAssignmentsClientWithBaseURI(sdkApi sdkEnv.Api) (*ConfigurationProfileAssignmentsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "configurationprofileassignments", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ConfigurationProfileAssignmentsClient: %+v", err)
	}

	return &ConfigurationProfileAssignmentsClient{
		Client: client,
	}, nil
}
