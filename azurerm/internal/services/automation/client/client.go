package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/automation/mgmt/2018-06-30-preview/automation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	AccountClient               *automation.AccountClient
	AgentRegistrationInfoClient *automation.AgentRegistrationInformationClient
	CertificateClient           *automation.CertificateClient
	ConnectionClient            *automation.ConnectionClient
	ConnectionTypeClient        *automation.ConnectionTypeClient
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

	connectionClient := automation.NewConnectionClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&connectionClient.Client, o.ResourceManagerAuthorizer)

	connectionTypeClient := automation.NewConnectionTypeClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&connectionTypeClient.Client, o.ResourceManagerAuthorizer)

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
		ConnectionClient:            &connectionClient,
		ConnectionTypeClient:        &connectionTypeClient,
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
