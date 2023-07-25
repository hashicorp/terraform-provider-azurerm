// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

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

func NewClient(o *common.ClientOptions) (*Client, error) {
	integrationAccountClient, err := integrationaccounts.NewIntegrationAccountsClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(integrationAccountClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building IntegrationAccountClient client: %+v", err)
	}

	integrationAccountAgreementClient, err := integrationaccountagreements.NewIntegrationAccountAgreementsClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(integrationAccountAgreementClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building IntegrationAccountAgreementClient client: %+v", err)
	}

	integrationAccountAssemblyClient, err := integrationaccountassemblies.NewIntegrationAccountAssembliesClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(integrationAccountAssemblyClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building IntegrationAccountAssemblyClient client: %+v", err)
	}

	integrationAccountBatchConfigurationClient, err := integrationaccountbatchconfigurations.NewIntegrationAccountBatchConfigurationsClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(integrationAccountBatchConfigurationClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building IntegrationAccountBatchConfigurationClient client: %+v", err)
	}

	integrationAccountCertificateClient, err := integrationaccountcertificates.NewIntegrationAccountCertificatesClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(integrationAccountCertificateClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building IntegrationAccountCertificateClient client: %+v", err)
	}

	integrationAccountMapClient, err := integrationaccountmaps.NewIntegrationAccountMapsClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(integrationAccountMapClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building IntegrationAccountMapClient client: %+v", err)
	}

	integrationAccountPartnerClient, err := integrationaccountpartners.NewIntegrationAccountPartnersClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(integrationAccountPartnerClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building IntegrationAccountPartnerClient client: %+v", err)
	}

	integrationAccountSchemaClient, err := integrationaccountschemas.NewIntegrationAccountSchemasClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(integrationAccountSchemaClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building IntegrationAccountSchemaClient client: %+v", err)
	}

	integrationAccountSessionClient, err := integrationaccountsessions.NewIntegrationAccountSessionsClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(integrationAccountSessionClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building IntegrationAccountSessionClient client: %+v", err)
	}

	integrationServiceEnvironmentClient, err := integrationserviceenvironments.NewIntegrationServiceEnvironmentsClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(integrationServiceEnvironmentClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building IntegrationServiceEnvironmentClient client: %+v", err)
	}

	workflowClient, err := workflows.NewWorkflowsClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(workflowClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building WorkflowClient client: %+v", err)
	}

	triggersClient, err := workflowtriggers.NewWorkflowTriggersClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(triggersClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building TriggersClient client: %+v", err)
	}

	return &Client{
		IntegrationAccountClient:                   integrationAccountClient,
		IntegrationAccountAgreementClient:          integrationAccountAgreementClient,
		IntegrationAccountAssemblyClient:           integrationAccountAssemblyClient,
		IntegrationAccountBatchConfigurationClient: integrationAccountBatchConfigurationClient,
		IntegrationAccountCertificateClient:        integrationAccountCertificateClient,
		IntegrationAccountMapClient:                integrationAccountMapClient,
		IntegrationAccountPartnerClient:            integrationAccountPartnerClient,
		IntegrationAccountSchemaClient:             integrationAccountSchemaClient,
		IntegrationAccountSessionClient:            integrationAccountSessionClient,
		IntegrationServiceEnvironmentClient:        integrationServiceEnvironmentClient,
		WorkflowClient:                             workflowClient,
		TriggersClient:                             triggersClient,
	}, nil
}
