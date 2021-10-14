package client

import (
	"github.com/Azure/azure-sdk-for-go/services/logic/mgmt/2019-05-01/logic"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	IntegrationAccountClient                   *logic.IntegrationAccountsClient
	IntegrationAccountAgreementClient          *logic.IntegrationAccountAgreementsClient
	IntegrationAccountAssemblyClient           *logic.IntegrationAccountAssembliesClient
	IntegrationAccountBatchConfigurationClient *logic.IntegrationAccountBatchConfigurationsClient
	IntegrationAccountCertificateClient        *logic.IntegrationAccountCertificatesClient
	IntegrationAccountMapClient                *logic.IntegrationAccountMapsClient
	IntegrationAccountPartnerClient            *logic.IntegrationAccountPartnersClient
	IntegrationAccountSchemaClient             *logic.IntegrationAccountSchemasClient
	IntegrationAccountSessionClient            *logic.IntegrationAccountSessionsClient
	IntegrationServiceEnvironmentClient        *logic.IntegrationServiceEnvironmentsClient
	WorkflowClient                             *logic.WorkflowsClient
	TriggersClient                             *logic.WorkflowTriggersClient
}

func NewClient(o *common.ClientOptions) *Client {
	integrationAccountClient := logic.NewIntegrationAccountsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&integrationAccountClient.Client, o.ResourceManagerAuthorizer)

	integrationAccountAgreementClient := logic.NewIntegrationAccountAgreementsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&integrationAccountAgreementClient.Client, o.ResourceManagerAuthorizer)

	integrationAccountAssemblyClient := logic.NewIntegrationAccountAssembliesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&integrationAccountAssemblyClient.Client, o.ResourceManagerAuthorizer)

	integrationAccountBatchConfigurationClient := logic.NewIntegrationAccountBatchConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&integrationAccountBatchConfigurationClient.Client, o.ResourceManagerAuthorizer)

	integrationAccountCertificateClient := logic.NewIntegrationAccountCertificatesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&integrationAccountCertificateClient.Client, o.ResourceManagerAuthorizer)

	integrationAccountMapClient := logic.NewIntegrationAccountMapsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&integrationAccountMapClient.Client, o.ResourceManagerAuthorizer)

	integrationAccountPartnerClient := logic.NewIntegrationAccountPartnersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&integrationAccountPartnerClient.Client, o.ResourceManagerAuthorizer)

	integrationAccountSchemaClient := logic.NewIntegrationAccountSchemasClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&integrationAccountSchemaClient.Client, o.ResourceManagerAuthorizer)

	integrationAccountSessionClient := logic.NewIntegrationAccountSessionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&integrationAccountSessionClient.Client, o.ResourceManagerAuthorizer)

	integrationServiceEnvironmentClient := logic.NewIntegrationServiceEnvironmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&integrationServiceEnvironmentClient.Client, o.ResourceManagerAuthorizer)

	workflowClient := logic.NewWorkflowsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&workflowClient.Client, o.ResourceManagerAuthorizer)

	triggersClient := logic.NewWorkflowTriggersClient(o.SubscriptionId)
	o.ConfigureClient(&triggersClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		IntegrationAccountClient:                   &integrationAccountClient,
		IntegrationAccountAgreementClient:          &integrationAccountAgreementClient,
		IntegrationAccountAssemblyClient:           &integrationAccountAssemblyClient,
		IntegrationAccountBatchConfigurationClient: &integrationAccountBatchConfigurationClient,
		IntegrationAccountCertificateClient:        &integrationAccountCertificateClient,
		IntegrationAccountMapClient:                &integrationAccountMapClient,
		IntegrationAccountPartnerClient:            &integrationAccountPartnerClient,
		IntegrationAccountSchemaClient:             &integrationAccountSchemaClient,
		IntegrationAccountSessionClient:            &integrationAccountSessionClient,
		IntegrationServiceEnvironmentClient:        &integrationServiceEnvironmentClient,
		WorkflowClient:                             &workflowClient,
		TriggersClient:                             &triggersClient,
	}
}
