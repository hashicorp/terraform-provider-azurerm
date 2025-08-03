// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2021-06-01-preview/policy" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/guestconfiguration/2020-06-25/guestconfigurationassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/policyinsights/2021-10-01/remediations"
	assignments "github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-06-01/policyassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2025-01-01/policysetdefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AssignmentsClient                   *assignments.PolicyAssignmentsClient
	GuestConfigurationAssignmentsClient *guestconfigurationassignments.GuestConfigurationAssignmentsClient
	PolicySetDefinitionsClient          *policysetdefinitions.PolicySetDefinitionsClient
	RemediationsClient                  *remediations.RemediationsClient

	// TODO: Migrate these clients to `go-azure-sdk`
	DefinitionsClient    *policy.DefinitionsClient
	ExemptionsClient     *policy.ExemptionsClient
	SetDefinitionsClient *policy.SetDefinitionsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	assignmentsClient, err := assignments.NewPolicyAssignmentsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building PolicyAssignments client: %+v", err)
	}
	o.Configure(assignmentsClient.Client, o.Authorizers.ResourceManager)

	guestConfigurationAssignmentsClient, err := guestconfigurationassignments.NewGuestConfigurationAssignmentsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Guest Configuration Assignments Client:  %+v", err)
	}
	o.Configure(guestConfigurationAssignmentsClient.Client, o.Authorizers.ResourceManager)

	policySetDefinitionsClient, err := policysetdefinitions.NewPolicySetDefinitionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Policy Set Definitions client: %+v", err)
	}
	o.Configure(policySetDefinitionsClient.Client, o.Authorizers.ResourceManager)

	remediationsClient, err := remediations.NewRemediationsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Remediations client: %+v", err)
	}
	o.Configure(remediationsClient.Client, o.Authorizers.ResourceManager)

	// Track 1
	definitionsClient := policy.NewDefinitionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&definitionsClient.Client, o.ResourceManagerAuthorizer)

	exemptionsClient := policy.NewExemptionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&exemptionsClient.Client, o.ResourceManagerAuthorizer)

	setDefinitionsClient := policy.NewSetDefinitionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&setDefinitionsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AssignmentsClient:                   assignmentsClient,
		GuestConfigurationAssignmentsClient: guestConfigurationAssignmentsClient,
		PolicySetDefinitionsClient:          policySetDefinitionsClient,
		RemediationsClient:                  remediationsClient,

		// Track 1
		DefinitionsClient:    &definitionsClient,
		ExemptionsClient:     &exemptionsClient,
		SetDefinitionsClient: &setDefinitionsClient,
	}, nil
}
