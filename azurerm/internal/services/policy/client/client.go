package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/policyinsights/mgmt/2019-10-01-preview/policyinsights"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-09-01/policy"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	AssignmentsClient    *policy.AssignmentsClient
	DefinitionsClient    *policy.DefinitionsClient
	SetDefinitionsClient *policy.SetDefinitionsClient
	RemediationsClient   *policyinsights.RemediationsClient
}

func NewClient(o *common.ClientOptions) *Client {
	assignmentsClient := policy.NewAssignmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&assignmentsClient.Client, o.ResourceManagerAuthorizer)

	definitionsClient := policy.NewDefinitionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&definitionsClient.Client, o.ResourceManagerAuthorizer)

	setDefinitionsClient := policy.NewSetDefinitionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&setDefinitionsClient.Client, o.ResourceManagerAuthorizer)

	remediationsClient := policyinsights.NewRemediationsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&remediationsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AssignmentsClient:    &assignmentsClient,
		DefinitionsClient:    &definitionsClient,
		SetDefinitionsClient: &setDefinitionsClient,
		RemediationsClient:   &remediationsClient,
	}
}
