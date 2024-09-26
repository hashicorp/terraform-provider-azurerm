// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofileassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofilehciassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofilehcrpassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ConfigurationProfilesClient                     *configurationprofiles.ConfigurationProfilesClient
	ConfigurationProfileArcMachineAssignmentsClient *configurationprofilehcrpassignments.ConfigurationProfileHCRPAssignmentsClient
	ConfigurationProfileHCIAssignmentsClient        *configurationprofilehciassignments.ConfigurationProfileHCIAssignmentsClient
	ConfigurationProfileVMAssignmentsClient         *configurationprofileassignments.ConfigurationProfileAssignmentsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	configurationProfilesClient, err := configurationprofiles.NewConfigurationProfilesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ConfigurationProfiles client: %+v", err)
	}
	o.Configure(configurationProfilesClient.Client, o.Authorizers.ResourceManager)

	configurationProfileArcMachineAssignmentsClient, err := configurationprofilehcrpassignments.NewConfigurationProfileHCRPAssignmentsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ConfigurationProfilesHCIAssignments client: %+v", err)
	}
	o.Configure(configurationProfileArcMachineAssignmentsClient.Client, o.Authorizers.ResourceManager)

	configurationProfileHCIAssignmentsClient, err := configurationprofilehciassignments.NewConfigurationProfileHCIAssignmentsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ConfigurationProfilesHCIAssignments client: %+v", err)
	}
	o.Configure(configurationProfileHCIAssignmentsClient.Client, o.Authorizers.ResourceManager)

	configurationProfileVMAssignmentsClient, err := configurationprofileassignments.NewConfigurationProfileAssignmentsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ConfigurationProfilesHCIAssignments client: %+v", err)
	}
	o.Configure(configurationProfileVMAssignmentsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		ConfigurationProfilesClient:                     configurationProfilesClient,
		ConfigurationProfileArcMachineAssignmentsClient: configurationProfileArcMachineAssignmentsClient,
		ConfigurationProfileHCIAssignmentsClient:        configurationProfileHCIAssignmentsClient,
		ConfigurationProfileVMAssignmentsClient:         configurationProfileVMAssignmentsClient,
	}, nil
}
