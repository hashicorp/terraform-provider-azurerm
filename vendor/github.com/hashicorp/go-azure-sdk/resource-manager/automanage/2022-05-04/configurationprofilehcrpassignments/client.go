package configurationprofilehcrpassignments

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigurationProfileHCRPAssignmentsClient struct {
	Client *resourcemanager.Client
}

func NewConfigurationProfileHCRPAssignmentsClientWithBaseURI(sdkApi sdkEnv.Api) (*ConfigurationProfileHCRPAssignmentsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "configurationprofilehcrpassignments", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ConfigurationProfileHCRPAssignmentsClient: %+v", err)
	}

	return &ConfigurationProfileHCRPAssignmentsClient{
		Client: client,
	}, nil
}
