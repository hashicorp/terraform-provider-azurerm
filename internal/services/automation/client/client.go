package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/automation/mgmt/2020-01-13-preview/automation"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2019-06-01/runbook"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2021-06-22/automationaccount"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2021-06-22/hybridrunbookworker"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2021-06-22/hybridrunbookworkergroup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AccountClient               *automationaccount.AutomationAccountClient
	AgentRegistrationInfoClient *automation.AgentRegistrationInformationClient
	CertificateClient           *automation.CertificateClient
	ConnectionClient            *automation.ConnectionClient
	ConnectionTypeClient        *automation.ConnectionTypeClient
	CredentialClient            *automation.CredentialClient
	DscConfigurationClient      *automation.DscConfigurationClient
	DscNodeConfigurationClient  *automation.DscNodeConfigurationClient
	JobScheduleClient           *automation.JobScheduleClient
	ModuleClient                *automation.ModuleClient
	RunbookClient               *runbook.RunbookClient
	RunbookClientHack           *automation.RunbookClient
	RunbookDraftClient          *automation.RunbookDraftClient
	RunBookWgClient             *hybridrunbookworkergroup.HybridRunbookWorkerGroupClient
	RunbookWorkerClient         *hybridrunbookworker.HybridRunbookWorkerClient
	ScheduleClient              *automation.ScheduleClient
	SoftwareUpdateConfigClient  *automation.SoftwareUpdateConfigurationsClient
	SourceControlClient         *automation.SourceControlClient
	VariableClient              *automation.VariableClient
	WatcherClient               *automation.WatcherClient
	WebhookClient               *automation.WebhookClient
}

func NewClient(o *common.ClientOptions) *Client {
	accountClient := automationaccount.NewAutomationAccountClientWithBaseURI(o.ResourceManagerEndpoint)
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

	runbookClient2 := automation.NewRunbookClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&runbookClient2.Client, o.ResourceManagerAuthorizer)

	runbookClient := runbook.NewRunbookClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&runbookClient.Client, o.ResourceManagerAuthorizer)

	runbookDraftClient := automation.NewRunbookDraftClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&runbookDraftClient.Client, o.ResourceManagerAuthorizer)

	runbookWgClient := hybridrunbookworkergroup.NewHybridRunbookWorkerGroupClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&runbookWgClient.Client, o.ResourceManagerAuthorizer)

	runbookWorkerClient := hybridrunbookworker.NewHybridRunbookWorkerClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&runbookWorkerClient.Client, o.ResourceManagerAuthorizer)

	sourceCtlClient := automation.NewSourceControlClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&sourceCtlClient.Client, o.ResourceManagerAuthorizer)

	scheduleClient := automation.NewScheduleClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&scheduleClient.Client, o.ResourceManagerAuthorizer)

	softUpClient := automation.NewSoftwareUpdateConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&softUpClient.Client, o.ResourceManagerAuthorizer)

	variableClient := automation.NewVariableClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&variableClient.Client, o.ResourceManagerAuthorizer)

	watcherClient := automation.NewWatcherClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&watcherClient.Client, o.ResourceManagerAuthorizer)

	webhookClient := automation.NewWebhookClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&webhookClient.Client, o.ResourceManagerAuthorizer)

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
		RunbookClientHack:           &runbookClient2,
		RunbookDraftClient:          &runbookDraftClient,
		RunBookWgClient:             &runbookWgClient,
		RunbookWorkerClient:         &runbookWorkerClient,
		ScheduleClient:              &scheduleClient,
		SoftwareUpdateConfigClient:  &softUpClient,
		SourceControlClient:         &sourceCtlClient,
		VariableClient:              &variableClient,
		WatcherClient:               &watcherClient,
		WebhookClient:               &webhookClient,
	}
}
