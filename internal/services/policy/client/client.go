package client

import (
	"github.com/Azure/azure-sdk-for-go/services/guestconfiguration/mgmt/2020-06-25/guestconfiguration"
	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2021-06-01-preview/policy"
	policyPreview "github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2021-06-01-preview/policy"
	"github.com/hashicorp/go-azure-sdk/resource-manager/policyinsights/2021-10-01/policyinsights"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AssignmentsClient                   *policy.AssignmentsClient
	DefinitionsClient                   *policy.DefinitionsClient
	ExemptionsClient                    *policyPreview.ExemptionsClient
	SetDefinitionsClient                *policy.SetDefinitionsClient
	PolicyInsightsClient                *policyinsights.PolicyInsightsClient
	GuestConfigurationAssignmentsClient *guestconfiguration.AssignmentsClient
}

func NewClient(o *common.ClientOptions) *Client {
	assignmentsClient := policy.NewAssignmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&assignmentsClient.Client, o.ResourceManagerAuthorizer)

	definitionsClient := policy.NewDefinitionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&definitionsClient.Client, o.ResourceManagerAuthorizer)

	exemptionsClient := policyPreview.NewExemptionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&exemptionsClient.Client, o.ResourceManagerAuthorizer)

	setDefinitionsClient := policy.NewSetDefinitionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&setDefinitionsClient.Client, o.ResourceManagerAuthorizer)

	policyInsightsClient := policyinsights.NewPolicyInsightsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&policyInsightsClient.Client, o.ResourceManagerAuthorizer)

	guestConfigurationAssignmentsClient := guestconfiguration.NewAssignmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&guestConfigurationAssignmentsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AssignmentsClient:                   &assignmentsClient,
		DefinitionsClient:                   &definitionsClient,
		ExemptionsClient:                    &exemptionsClient,
		SetDefinitionsClient:                &setDefinitionsClient,
		PolicyInsightsClient:                &policyInsightsClient,
		GuestConfigurationAssignmentsClient: &guestConfigurationAssignmentsClient,
	}
}
