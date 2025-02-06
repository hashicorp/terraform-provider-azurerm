// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/managedservices/2022-10-01/registrationassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedservices/2022-10-01/registrationdefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AssignmentsClient *registrationassignments.RegistrationAssignmentsClient
	DefinitionsClient *registrationdefinitions.RegistrationDefinitionsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	assignmentsClient, err := registrationassignments.NewRegistrationAssignmentsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building RegistrationAssignments Client: %+v", err)
	}
	o.Configure(assignmentsClient.Client, o.Authorizers.ResourceManager)

	definitionsClient, err := registrationdefinitions.NewRegistrationDefinitionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building RegistrationDefinitions Client: %+v", err)
	}
	o.Configure(definitionsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		AssignmentsClient: assignmentsClient,
		DefinitionsClient: definitionsClient,
	}, nil
}
