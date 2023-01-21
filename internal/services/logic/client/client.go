package client

import (
	"github.com/Azure/azure-sdk-for-go/services/logic/mgmt/2019-05-01/logic"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountagreements"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountassemblies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountbatchconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountcertificates"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountmaps"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountpartners"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountschemas"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	IntegrationAccountClient                   *integrationaccounts.IntegrationAccountsClient
	IntegrationAccountAgreementClient          *integrationaccountagreements.IntegrationAccountAgreementsClient
	IntegrationAccountAssemblyClient           *integrationaccountassemblies.IntegrationAccountAssembliesClient
	IntegrationAccountBatchConfigurationClient *integrationaccountbatchconfigurations.IntegrationAccountBatchConfigurationsClient
	IntegrationAccountCertificateClient        *integrationaccountcertificates.IntegrationAccountCertificatesClient
	IntegrationAccountMapClient                *integrationaccountmaps.IntegrationAccountMapsClient
	IntegrationAccountPartnerClient            *integrationaccountpartners.IntegrationAccountPartnersClient
	IntegrationAccountSchemaClient             *integrationaccountschemas.IntegrationAccountSchemasClient
	IntegrationAccountSessionClient            *logic.IntegrationAccountSessionsClient
	IntegrationServiceEnvironmentClient        *logic.IntegrationServiceEnvironmentsClient
	WorkflowClient                             *logic.WorkflowsClient
	TriggersClient                             *logic.WorkflowTriggersClient
}

func NewClient(o *common.ClientOptions) *Client {
	integrationAccountClient := integrationaccounts.NewIntegrationAccountsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&integrationAccountClient.Client, o.ResourceManagerAuthorizer)

	integrationAccountAgreementClient := integrationaccountagreements.NewIntegrationAccountAgreementsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&integrationAccountAgreementClient.Client, o.ResourceManagerAuthorizer)

	integrationAccountAssemblyClient := integrationaccountassemblies.NewIntegrationAccountAssembliesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&integrationAccountAssemblyClient.Client, o.ResourceManagerAuthorizer)

	integrationAccountBatchConfigurationClient := integrationaccountbatchconfigurations.NewIntegrationAccountBatchConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&integrationAccountBatchConfigurationClient.Client, o.ResourceManagerAuthorizer)

	integrationAccountCertificateClient := integrationaccountcertificates.NewIntegrationAccountCertificatesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&integrationAccountCertificateClient.Client, o.ResourceManagerAuthorizer)

	integrationAccountMapClient := integrationaccountmaps.NewIntegrationAccountMapsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&integrationAccountMapClient.Client, o.ResourceManagerAuthorizer)

	integrationAccountPartnerClient := integrationaccountpartners.NewIntegrationAccountPartnersClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&integrationAccountPartnerClient.Client, o.ResourceManagerAuthorizer)

	integrationAccountSchemaClient := integrationaccountschemas.NewIntegrationAccountSchemasClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&integrationAccountSchemaClient.Client, o.ResourceManagerAuthorizer)

	integrationAccountSessionClient := logic.NewIntegrationAccountSessionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&integrationAccountSessionClient.Client, o.ResourceManagerAuthorizer)

	integrationServiceEnvironmentClient := logic.NewIntegrationServiceEnvironmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&integrationServiceEnvironmentClient.Client, o.ResourceManagerAuthorizer)

	workflowClient := logic.NewWorkflowsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&workflowClient.Client, o.ResourceManagerAuthorizer)

	triggersClient := logic.NewWorkflowTriggersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
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
