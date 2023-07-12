// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountagreements"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountassemblies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountbatchconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountcertificates"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountmaps"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountpartners"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountschemas"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountsessions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationserviceenvironments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflows"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflowtriggers"
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
	IntegrationAccountSessionClient            *integrationaccountsessions.IntegrationAccountSessionsClient
	IntegrationServiceEnvironmentClient        *integrationserviceenvironments.IntegrationServiceEnvironmentsClient
	WorkflowClient                             *workflows.WorkflowsClient
	TriggersClient                             *workflowtriggers.WorkflowTriggersClient
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

	integrationAccountSessionClient := integrationaccountsessions.NewIntegrationAccountSessionsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&integrationAccountSessionClient.Client, o.ResourceManagerAuthorizer)

	integrationServiceEnvironmentClient := integrationserviceenvironments.NewIntegrationServiceEnvironmentsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&integrationServiceEnvironmentClient.Client, o.ResourceManagerAuthorizer)

	workflowClient := workflows.NewWorkflowsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&workflowClient.Client, o.ResourceManagerAuthorizer)

	triggersClient := workflowtriggers.NewWorkflowTriggersClientWithBaseURI(o.ResourceManagerEndpoint)
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
