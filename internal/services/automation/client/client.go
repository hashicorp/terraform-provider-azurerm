// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/automation/mgmt/2020-01-13-preview/automation" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2015-10-31/webhook"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2019-06-01/softwareupdateconfiguration"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2020-01-13-preview/watcher"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/dscnodeconfiguration"

	// hybridrunbookworkergroup v2022-08-08 issue: https://github.com/Azure/azure-rest-api-specs/issues/24740
	runbook2 "github.com/hashicorp/go-azure-sdk/resource-manager/automation/2019-06-01/runbook"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2021-06-22/hybridrunbookworkergroup"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/automationaccount"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/certificate"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/connection"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/connectiontype"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/credential"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/dscconfiguration"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/hybridrunbookworker"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/jobschedule"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/module"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/runbook"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/schedule"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/sourcecontrol"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/variable"
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
	// use new sdk once https://github.com/hashicorp/pandora/issues/2754 fixed
	RunbookClientHack *runbook2.RunbookClient
	// port to new sdk issue once https://github.com/hashicorp/pandora/issues/2753 fixed
	RunbookDraftClient         *automation.RunbookDraftClient
	RunBookWgClient            *hybridrunbookworkergroup.HybridRunbookWorkerGroupClient
	RunbookWorkerClient        *hybridrunbookworker.HybridRunbookWorkerClient
	ScheduleClient             *schedule.ScheduleClient
	SoftwareUpdateConfigClient *softwareupdateconfiguration.SoftwareUpdateConfigurationClient
	SourceControlClient        *sourcecontrol.SourceControlClient
	VariableClient             *variable.VariableClient
	WatcherClient              *watcher.WatcherClient
	WebhookClient              *webhook.WebhookClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	agentRegistrationInfoClient := automation.NewAgentRegistrationInformationClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&agentRegistrationInfoClient.Client, o.ResourceManagerAuthorizer)

	runbookWgClient := hybridrunbookworkergroup.NewHybridRunbookWorkerGroupClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&runbookWgClient.Client, o.ResourceManagerAuthorizer)

	softUpClient := softwareupdateconfiguration.NewSoftwareUpdateConfigurationClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&softUpClient.Client, o.ResourceManagerAuthorizer)

	watcherClient := watcher.NewWatcherClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&watcherClient.Client, o.ResourceManagerAuthorizer)

	webhookClient := webhook.NewWebhookClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&webhookClient.Client, o.ResourceManagerAuthorizer)

	accountClient, err := automationaccount.NewAutomationAccountClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("build accountClient: %+v", err)
	}
	o.Configure(accountClient.Client, o.Authorizers.ResourceManager)

	certificateClient, err := certificate.NewCertificateClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("build certificateClient: %+v", err)
	}
	o.Configure(certificateClient.Client, o.Authorizers.ResourceManager)

	connectionClient, err := connection.NewConnectionClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("build connectionClient: %+v", err)
	}
	o.Configure(connectionClient.Client, o.Authorizers.ResourceManager)

	connectionTypeClient, err := connectiontype.NewConnectionTypeClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("build connectionTypeClient: %+v", err)
	}
	o.Configure(connectionTypeClient.Client, o.Authorizers.ResourceManager)

	credentialClient, err := credential.NewCredentialClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("build credentialClient: %+v", err)
	}
	o.Configure(credentialClient.Client, o.Authorizers.ResourceManager)

	dscConfigurationClient, err := dscconfiguration.NewDscConfigurationClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("build dscConfigurationClient: %+v", err)
	}
	o.Configure(dscConfigurationClient.Client, o.Authorizers.ResourceManager)

	dscNodeConfigurationClient, err := dscnodeconfiguration.NewDscNodeConfigurationClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building DscNodeConfiguration client: %+v", err)
	}
	o.Configure(dscNodeConfigurationClient.Client, o.Authorizers.ResourceManager)

	jobScheduleClient, err := jobschedule.NewJobScheduleClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("build jobScheduleClient: %+v", err)
	}
	o.Configure(jobScheduleClient.Client, o.Authorizers.ResourceManager)

	moduleClient, err := module.NewModuleClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("build moduleClient: %+v", err)
	}
	o.Configure(moduleClient.Client, o.Authorizers.ResourceManager)

	runbookClient2 := runbook2.NewRunbookClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&runbookClient2.Client, o.ResourceManagerAuthorizer)

	runbookClient, err := runbook.NewRunbookClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("build runbookClient: %+v", err)
	}
	o.Configure(runbookClient.Client, o.Authorizers.ResourceManager)

	runbookDraftClient := automation.NewRunbookDraftClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&runbookDraftClient.Client, o.ResourceManagerAuthorizer)

	// runbookDraftClient := runbookdraft.NewRunbookDraftClientWithBaseURI(o.ResourceManagerEndpoint)
	// o.ConfigureClient(&runbookDraftClient.Client, o.ResourceManagerAuthorizer)

	runbookWorkerClient, err := hybridrunbookworker.NewHybridRunbookWorkerClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("build runbookWorkerClient: %+v", err)
	}
	o.Configure(runbookWorkerClient.Client, o.Authorizers.ResourceManager)

	sourceCtlClient, err := sourcecontrol.NewSourceControlClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("build sourceCtlClient: %+v", err)
	}
	o.Configure(sourceCtlClient.Client, o.Authorizers.ResourceManager)

	scheduleClient, err := schedule.NewScheduleClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("build scheduleClient: %+v", err)
	}
	o.Configure(scheduleClient.Client, o.Authorizers.ResourceManager)

	variableClient, err := variable.NewVariableClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("build variableClient: %+v", err)
	}
	o.Configure(variableClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		AccountClient:               accountClient,
		AgentRegistrationInfoClient: &agentRegistrationInfoClient,
		CertificateClient:           certificateClient,
		ConnectionClient:            connectionClient,
		ConnectionTypeClient:        connectionTypeClient,
		CredentialClient:            credentialClient,
		DscConfigurationClient:      dscConfigurationClient,
		DscNodeConfigurationClient:  dscNodeConfigurationClient,
		JobScheduleClient:           jobScheduleClient,
		ModuleClient:                moduleClient,
		RunbookClient:               runbookClient,
		RunbookClientHack:           &runbookClient2,
		RunbookDraftClient:          &runbookDraftClient,
		RunBookWgClient:             &runbookWgClient,
		RunbookWorkerClient:         runbookWorkerClient,
		ScheduleClient:              scheduleClient,
		SoftwareUpdateConfigClient:  &softUpClient,
		SourceControlClient:         sourceCtlClient,
		VariableClient:              variableClient,
		WatcherClient:               &watcherClient,
		WebhookClient:               &webhookClient,
	}, nil
}
