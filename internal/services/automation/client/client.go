package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/automation/mgmt/2020-01-13-preview/automation" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2015-10-31/webhook"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2019-06-01/dscconfiguration"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2019-06-01/runbook"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2019-06-01/softwareupdateconfiguration"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2020-01-13-preview/certificate"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2020-01-13-preview/connection"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2020-01-13-preview/connectiontype"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2020-01-13-preview/credential"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2020-01-13-preview/dscnodeconfiguration"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2020-01-13-preview/jobschedule"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2020-01-13-preview/module"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2020-01-13-preview/schedule"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2020-01-13-preview/sourcecontrol"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2020-01-13-preview/variable"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2020-01-13-preview/watcher"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2021-06-22/automationaccount"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2021-06-22/hybridrunbookworker"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2021-06-22/hybridrunbookworkergroup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AccountClient               *automationaccount.AutomationAccountClient
	AgentRegistrationInfoClient *automation.AgentRegistrationInformationClient
	CertificateClient           *certificate.CertificateClient
	ConnectionClient            *connection.ConnectionClient
	ConnectionTypeClient        *connectiontype.ConnectionTypeClient
	CredentialClient            *credential.CredentialClient
	DscConfigurationClient      *dscconfiguration.DscConfigurationClient
	DscNodeConfigurationClient  *dscnodeconfiguration.DscNodeConfigurationClient
	JobScheduleClient           *jobschedule.JobScheduleClient
	ModuleClient                *module.ModuleClient
	RunbookClient               *runbook.RunbookClient
	RunbookClientHack           *automation.RunbookClient
	RunbookDraftClient          *automation.RunbookDraftClient
	RunBookWgClient             *hybridrunbookworkergroup.HybridRunbookWorkerGroupClient
	RunbookWorkerClient         *hybridrunbookworker.HybridRunbookWorkerClient
	ScheduleClient              *schedule.ScheduleClient
	SoftwareUpdateConfigClient  *softwareupdateconfiguration.SoftwareUpdateConfigurationClient
	SourceControlClient         *sourcecontrol.SourceControlClient
	VariableClient              *variable.VariableClient
	WatcherClient               *watcher.WatcherClient
	WebhookClient               *webhook.WebhookClient
}

func NewClient(o *common.ClientOptions) *Client {
	accountClient := automationaccount.NewAutomationAccountClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&accountClient.Client, o.ResourceManagerAuthorizer)

	agentRegistrationInfoClient := automation.NewAgentRegistrationInformationClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&agentRegistrationInfoClient.Client, o.ResourceManagerAuthorizer)

	certificateClient := certificate.NewCertificateClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&certificateClient.Client, o.ResourceManagerAuthorizer)

	connectionClient := connection.NewConnectionClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&connectionClient.Client, o.ResourceManagerAuthorizer)

	connectionTypeClient := connectiontype.NewConnectionTypeClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&connectionTypeClient.Client, o.ResourceManagerAuthorizer)

	credentialClient := credential.NewCredentialClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&credentialClient.Client, o.ResourceManagerAuthorizer)

	dscConfigurationClient := dscconfiguration.NewDscConfigurationClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&dscConfigurationClient.Client, o.ResourceManagerAuthorizer)

	dscNodeConfigurationClient := dscnodeconfiguration.NewDscNodeConfigurationClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&dscNodeConfigurationClient.Client, o.ResourceManagerAuthorizer)

	jobScheduleClient := jobschedule.NewJobScheduleClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&jobScheduleClient.Client, o.ResourceManagerAuthorizer)

	moduleClient := module.NewModuleClientWithBaseURI(o.ResourceManagerEndpoint)
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

	sourceCtlClient := sourcecontrol.NewSourceControlClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&sourceCtlClient.Client, o.ResourceManagerAuthorizer)

	scheduleClient := schedule.NewScheduleClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&scheduleClient.Client, o.ResourceManagerAuthorizer)

	softUpClient := softwareupdateconfiguration.NewSoftwareUpdateConfigurationClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&softUpClient.Client, o.ResourceManagerAuthorizer)

	variableClient := variable.NewVariableClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&variableClient.Client, o.ResourceManagerAuthorizer)

	watcherClient := watcher.NewWatcherClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&watcherClient.Client, o.ResourceManagerAuthorizer)

	webhookClient := webhook.NewWebhookClientWithBaseURI(o.ResourceManagerEndpoint)
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
