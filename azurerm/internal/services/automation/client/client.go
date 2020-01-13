package client

import (
	"github.com/Azure/azure-sdk-for-go/services/automation/mgmt/2015-10-31/automation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	AccountClient               *automation.AccountClient
	AgentRegistrationInfoClient *automation.AgentRegistrationInformationClient
	CertificateClient           *automation.CertificateClient
	CredentialClient            *automation.CredentialClient
	DscConfigurationClient      *automation.DscConfigurationClient
	DscNodeConfigurationClient  *automation.DscNodeConfigurationClient
	JobScheduleClient           *automation.JobScheduleClient
	ModuleClient                *automation.ModuleClient
	RunbookClient               *automation.RunbookClient
	RunbookDraftClient          *automation.RunbookDraftClient
	ScheduleClient              *automation.ScheduleClient
	VariableClient              *automation.VariableClient
}

func NewClient(o *common.ClientOptions) *Client {
	accountClient := automation.NewAccountClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&accountClient.Client, o.ResourceManagerAuthorizer)

	agentRegistrationInfoClient := automation.NewAgentRegistrationInformationClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&agentRegistrationInfoClient.Client, o.ResourceManagerAuthorizer)

	certificateClient := automation.NewCertificateClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&certificateClient.Client, o.ResourceManagerAuthorizer)

	credentialClient := automation.NewCredentialClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&credentialClient.Client, o.ResourceManagerAuthorizer)

	dscConfigurationClient := automation.NewDscConfigurationClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&dscConfigurationClient.Client, o.ResourceManagerAuthorizer)

	dscNodeConfigurationClient := automation.NewDscNodeConfigurationClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&dscNodeConfigurationClient.Client, o.ResourceManagerAuthorizer)

	jobScheduleClient := automation.NewJobScheduleClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&jobScheduleClient.Client, o.ResourceManagerAuthorizer)

	moduleClient := automation.NewModuleClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&moduleClient.Client, o.ResourceManagerAuthorizer)

	runbookClient := automation.NewRunbookClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&runbookClient.Client, o.ResourceManagerAuthorizer)

	runbookDraftClient := automation.NewRunbookDraftClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&runbookDraftClient.Client, o.ResourceManagerAuthorizer)

	scheduleClient := automation.NewScheduleClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&scheduleClient.Client, o.ResourceManagerAuthorizer)

	variableClient := automation.NewVariableClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&variableClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AccountClient:               &accountClient,
		AgentRegistrationInfoClient: &agentRegistrationInfoClient,
		CertificateClient:           &certificateClient,
		CredentialClient:            &credentialClient,
		DscConfigurationClient:      &dscConfigurationClient,
		DscNodeConfigurationClient:  &dscNodeConfigurationClient,
		JobScheduleClient:           &jobScheduleClient,
		ModuleClient:                &moduleClient,
		RunbookClient:               &runbookClient,
		RunbookDraftClient:          &runbookDraftClient,
		ScheduleClient:              &scheduleClient,
		VariableClient:              &variableClient,
	}
}
