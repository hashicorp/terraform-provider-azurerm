package v2024_10_23

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/activity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/agentregistrationinformation"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/automationaccount"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/certificate"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/connection"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/connectiontype"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/credential"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/deletedautomationaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/dscconfiguration"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/dscnode"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/dscnodeconfiguration"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/hybridrunbookworker"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/hybridrunbookworkergroup"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/job"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/jobschedule"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/jobstream"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/linkedworkspace"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/listallhybridrunbookworkergroupinautomationaccount"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/listdeletedrunbooks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/listkeys"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/module"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/nodecountinformation"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/nodereports"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/objectdatatypes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/operations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/packageresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/privateendpointconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/privatelinkresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/python2package"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/python3package"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/runbook"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/runbookdraft"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/runtimeenvironment"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/schedule"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/softwareupdateconfiguration"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/softwareupdateconfigurationmachinerun"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/softwareupdateconfigurationrun"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/sourcecontrol"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/sourcecontrolsyncjob"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/sourcecontrolsyncjobstreams"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/statistics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/testjob"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/testjobstream"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/typefields"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/usages"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/variable"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/watcher"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/webhook"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	Activity                                           *activity.ActivityClient
	AgentRegistrationInformation                       *agentregistrationinformation.AgentRegistrationInformationClient
	AutomationAccount                                  *automationaccount.AutomationAccountClient
	Certificate                                        *certificate.CertificateClient
	Connection                                         *connection.ConnectionClient
	ConnectionType                                     *connectiontype.ConnectionTypeClient
	Credential                                         *credential.CredentialClient
	DeletedAutomationAccounts                          *deletedautomationaccounts.DeletedAutomationAccountsClient
	DscConfiguration                                   *dscconfiguration.DscConfigurationClient
	DscNode                                            *dscnode.DscNodeClient
	DscNodeConfiguration                               *dscnodeconfiguration.DscNodeConfigurationClient
	HybridRunbookWorker                                *hybridrunbookworker.HybridRunbookWorkerClient
	HybridRunbookWorkerGroup                           *hybridrunbookworkergroup.HybridRunbookWorkerGroupClient
	Job                                                *job.JobClient
	JobSchedule                                        *jobschedule.JobScheduleClient
	JobStream                                          *jobstream.JobStreamClient
	LinkedWorkspace                                    *linkedworkspace.LinkedWorkspaceClient
	ListAllHybridRunbookWorkerGroupInAutomationAccount *listallhybridrunbookworkergroupinautomationaccount.ListAllHybridRunbookWorkerGroupInAutomationAccountClient
	ListDeletedRunbooks                                *listdeletedrunbooks.ListDeletedRunbooksClient
	ListKeys                                           *listkeys.ListKeysClient
	Module                                             *module.ModuleClient
	NodeCountInformation                               *nodecountinformation.NodeCountInformationClient
	NodeReports                                        *nodereports.NodeReportsClient
	ObjectDataTypes                                    *objectdatatypes.ObjectDataTypesClient
	Operations                                         *operations.OperationsClient
	PackageResource                                    *packageresource.PackageResourceClient
	PrivateEndpointConnections                         *privateendpointconnections.PrivateEndpointConnectionsClient
	PrivateLinkResources                               *privatelinkresources.PrivateLinkResourcesClient
	Python2Package                                     *python2package.Python2PackageClient
	Python3Package                                     *python3package.Python3PackageClient
	Runbook                                            *runbook.RunbookClient
	RunbookDraft                                       *runbookdraft.RunbookDraftClient
	RuntimeEnvironment                                 *runtimeenvironment.RuntimeEnvironmentClient
	Schedule                                           *schedule.ScheduleClient
	SoftwareUpdateConfiguration                        *softwareupdateconfiguration.SoftwareUpdateConfigurationClient
	SoftwareUpdateConfigurationMachineRun              *softwareupdateconfigurationmachinerun.SoftwareUpdateConfigurationMachineRunClient
	SoftwareUpdateConfigurationRun                     *softwareupdateconfigurationrun.SoftwareUpdateConfigurationRunClient
	SourceControl                                      *sourcecontrol.SourceControlClient
	SourceControlSyncJob                               *sourcecontrolsyncjob.SourceControlSyncJobClient
	SourceControlSyncJobStreams                        *sourcecontrolsyncjobstreams.SourceControlSyncJobStreamsClient
	Statistics                                         *statistics.StatisticsClient
	TestJob                                            *testjob.TestJobClient
	TestJobStream                                      *testjobstream.TestJobStreamClient
	TypeFields                                         *typefields.TypeFieldsClient
	Usages                                             *usages.UsagesClient
	Variable                                           *variable.VariableClient
	Watcher                                            *watcher.WatcherClient
	Webhook                                            *webhook.WebhookClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	activityClient, err := activity.NewActivityClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Activity client: %+v", err)
	}
	configureFunc(activityClient.Client)

	agentRegistrationInformationClient, err := agentregistrationinformation.NewAgentRegistrationInformationClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building AgentRegistrationInformation client: %+v", err)
	}
	configureFunc(agentRegistrationInformationClient.Client)

	automationAccountClient, err := automationaccount.NewAutomationAccountClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building AutomationAccount client: %+v", err)
	}
	configureFunc(automationAccountClient.Client)

	certificateClient, err := certificate.NewCertificateClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Certificate client: %+v", err)
	}
	configureFunc(certificateClient.Client)

	connectionClient, err := connection.NewConnectionClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Connection client: %+v", err)
	}
	configureFunc(connectionClient.Client)

	connectionTypeClient, err := connectiontype.NewConnectionTypeClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ConnectionType client: %+v", err)
	}
	configureFunc(connectionTypeClient.Client)

	credentialClient, err := credential.NewCredentialClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Credential client: %+v", err)
	}
	configureFunc(credentialClient.Client)

	deletedAutomationAccountsClient, err := deletedautomationaccounts.NewDeletedAutomationAccountsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building DeletedAutomationAccounts client: %+v", err)
	}
	configureFunc(deletedAutomationAccountsClient.Client)

	dscConfigurationClient, err := dscconfiguration.NewDscConfigurationClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building DscConfiguration client: %+v", err)
	}
	configureFunc(dscConfigurationClient.Client)

	dscNodeClient, err := dscnode.NewDscNodeClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building DscNode client: %+v", err)
	}
	configureFunc(dscNodeClient.Client)

	dscNodeConfigurationClient, err := dscnodeconfiguration.NewDscNodeConfigurationClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building DscNodeConfiguration client: %+v", err)
	}
	configureFunc(dscNodeConfigurationClient.Client)

	hybridRunbookWorkerClient, err := hybridrunbookworker.NewHybridRunbookWorkerClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building HybridRunbookWorker client: %+v", err)
	}
	configureFunc(hybridRunbookWorkerClient.Client)

	hybridRunbookWorkerGroupClient, err := hybridrunbookworkergroup.NewHybridRunbookWorkerGroupClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building HybridRunbookWorkerGroup client: %+v", err)
	}
	configureFunc(hybridRunbookWorkerGroupClient.Client)

	jobClient, err := job.NewJobClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Job client: %+v", err)
	}
	configureFunc(jobClient.Client)

	jobScheduleClient, err := jobschedule.NewJobScheduleClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building JobSchedule client: %+v", err)
	}
	configureFunc(jobScheduleClient.Client)

	jobStreamClient, err := jobstream.NewJobStreamClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building JobStream client: %+v", err)
	}
	configureFunc(jobStreamClient.Client)

	linkedWorkspaceClient, err := linkedworkspace.NewLinkedWorkspaceClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building LinkedWorkspace client: %+v", err)
	}
	configureFunc(linkedWorkspaceClient.Client)

	listAllHybridRunbookWorkerGroupInAutomationAccountClient, err := listallhybridrunbookworkergroupinautomationaccount.NewListAllHybridRunbookWorkerGroupInAutomationAccountClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ListAllHybridRunbookWorkerGroupInAutomationAccount client: %+v", err)
	}
	configureFunc(listAllHybridRunbookWorkerGroupInAutomationAccountClient.Client)

	listDeletedRunbooksClient, err := listdeletedrunbooks.NewListDeletedRunbooksClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ListDeletedRunbooks client: %+v", err)
	}
	configureFunc(listDeletedRunbooksClient.Client)

	listKeysClient, err := listkeys.NewListKeysClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ListKeys client: %+v", err)
	}
	configureFunc(listKeysClient.Client)

	moduleClient, err := module.NewModuleClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Module client: %+v", err)
	}
	configureFunc(moduleClient.Client)

	nodeCountInformationClient, err := nodecountinformation.NewNodeCountInformationClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building NodeCountInformation client: %+v", err)
	}
	configureFunc(nodeCountInformationClient.Client)

	nodeReportsClient, err := nodereports.NewNodeReportsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building NodeReports client: %+v", err)
	}
	configureFunc(nodeReportsClient.Client)

	objectDataTypesClient, err := objectdatatypes.NewObjectDataTypesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ObjectDataTypes client: %+v", err)
	}
	configureFunc(objectDataTypesClient.Client)

	operationsClient, err := operations.NewOperationsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Operations client: %+v", err)
	}
	configureFunc(operationsClient.Client)

	packageResourceClient, err := packageresource.NewPackageResourceClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building PackageResource client: %+v", err)
	}
	configureFunc(packageResourceClient.Client)

	privateEndpointConnectionsClient, err := privateendpointconnections.NewPrivateEndpointConnectionsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building PrivateEndpointConnections client: %+v", err)
	}
	configureFunc(privateEndpointConnectionsClient.Client)

	privateLinkResourcesClient, err := privatelinkresources.NewPrivateLinkResourcesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building PrivateLinkResources client: %+v", err)
	}
	configureFunc(privateLinkResourcesClient.Client)

	python2PackageClient, err := python2package.NewPython2PackageClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Python2Package client: %+v", err)
	}
	configureFunc(python2PackageClient.Client)

	python3PackageClient, err := python3package.NewPython3PackageClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Python3Package client: %+v", err)
	}
	configureFunc(python3PackageClient.Client)

	runbookClient, err := runbook.NewRunbookClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Runbook client: %+v", err)
	}
	configureFunc(runbookClient.Client)

	runbookDraftClient, err := runbookdraft.NewRunbookDraftClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building RunbookDraft client: %+v", err)
	}
	configureFunc(runbookDraftClient.Client)

	runtimeEnvironmentClient, err := runtimeenvironment.NewRuntimeEnvironmentClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building RuntimeEnvironment client: %+v", err)
	}
	configureFunc(runtimeEnvironmentClient.Client)

	scheduleClient, err := schedule.NewScheduleClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Schedule client: %+v", err)
	}
	configureFunc(scheduleClient.Client)

	softwareUpdateConfigurationClient, err := softwareupdateconfiguration.NewSoftwareUpdateConfigurationClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building SoftwareUpdateConfiguration client: %+v", err)
	}
	configureFunc(softwareUpdateConfigurationClient.Client)

	softwareUpdateConfigurationMachineRunClient, err := softwareupdateconfigurationmachinerun.NewSoftwareUpdateConfigurationMachineRunClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building SoftwareUpdateConfigurationMachineRun client: %+v", err)
	}
	configureFunc(softwareUpdateConfigurationMachineRunClient.Client)

	softwareUpdateConfigurationRunClient, err := softwareupdateconfigurationrun.NewSoftwareUpdateConfigurationRunClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building SoftwareUpdateConfigurationRun client: %+v", err)
	}
	configureFunc(softwareUpdateConfigurationRunClient.Client)

	sourceControlClient, err := sourcecontrol.NewSourceControlClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building SourceControl client: %+v", err)
	}
	configureFunc(sourceControlClient.Client)

	sourceControlSyncJobClient, err := sourcecontrolsyncjob.NewSourceControlSyncJobClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building SourceControlSyncJob client: %+v", err)
	}
	configureFunc(sourceControlSyncJobClient.Client)

	sourceControlSyncJobStreamsClient, err := sourcecontrolsyncjobstreams.NewSourceControlSyncJobStreamsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building SourceControlSyncJobStreams client: %+v", err)
	}
	configureFunc(sourceControlSyncJobStreamsClient.Client)

	statisticsClient, err := statistics.NewStatisticsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Statistics client: %+v", err)
	}
	configureFunc(statisticsClient.Client)

	testJobClient, err := testjob.NewTestJobClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building TestJob client: %+v", err)
	}
	configureFunc(testJobClient.Client)

	testJobStreamClient, err := testjobstream.NewTestJobStreamClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building TestJobStream client: %+v", err)
	}
	configureFunc(testJobStreamClient.Client)

	typeFieldsClient, err := typefields.NewTypeFieldsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building TypeFields client: %+v", err)
	}
	configureFunc(typeFieldsClient.Client)

	usagesClient, err := usages.NewUsagesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Usages client: %+v", err)
	}
	configureFunc(usagesClient.Client)

	variableClient, err := variable.NewVariableClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Variable client: %+v", err)
	}
	configureFunc(variableClient.Client)

	watcherClient, err := watcher.NewWatcherClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Watcher client: %+v", err)
	}
	configureFunc(watcherClient.Client)

	webhookClient, err := webhook.NewWebhookClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Webhook client: %+v", err)
	}
	configureFunc(webhookClient.Client)

	return &Client{
		Activity:                     activityClient,
		AgentRegistrationInformation: agentRegistrationInformationClient,
		AutomationAccount:            automationAccountClient,
		Certificate:                  certificateClient,
		Connection:                   connectionClient,
		ConnectionType:               connectionTypeClient,
		Credential:                   credentialClient,
		DeletedAutomationAccounts:    deletedAutomationAccountsClient,
		DscConfiguration:             dscConfigurationClient,
		DscNode:                      dscNodeClient,
		DscNodeConfiguration:         dscNodeConfigurationClient,
		HybridRunbookWorker:          hybridRunbookWorkerClient,
		HybridRunbookWorkerGroup:     hybridRunbookWorkerGroupClient,
		Job:                          jobClient,
		JobSchedule:                  jobScheduleClient,
		JobStream:                    jobStreamClient,
		LinkedWorkspace:              linkedWorkspaceClient,
		ListAllHybridRunbookWorkerGroupInAutomationAccount: listAllHybridRunbookWorkerGroupInAutomationAccountClient,
		ListDeletedRunbooks:                   listDeletedRunbooksClient,
		ListKeys:                              listKeysClient,
		Module:                                moduleClient,
		NodeCountInformation:                  nodeCountInformationClient,
		NodeReports:                           nodeReportsClient,
		ObjectDataTypes:                       objectDataTypesClient,
		Operations:                            operationsClient,
		PackageResource:                       packageResourceClient,
		PrivateEndpointConnections:            privateEndpointConnectionsClient,
		PrivateLinkResources:                  privateLinkResourcesClient,
		Python2Package:                        python2PackageClient,
		Python3Package:                        python3PackageClient,
		Runbook:                               runbookClient,
		RunbookDraft:                          runbookDraftClient,
		RuntimeEnvironment:                    runtimeEnvironmentClient,
		Schedule:                              scheduleClient,
		SoftwareUpdateConfiguration:           softwareUpdateConfigurationClient,
		SoftwareUpdateConfigurationMachineRun: softwareUpdateConfigurationMachineRunClient,
		SoftwareUpdateConfigurationRun:        softwareUpdateConfigurationRunClient,
		SourceControl:                         sourceControlClient,
		SourceControlSyncJob:                  sourceControlSyncJobClient,
		SourceControlSyncJobStreams:           sourceControlSyncJobStreamsClient,
		Statistics:                            statisticsClient,
		TestJob:                               testJobClient,
		TestJobStream:                         testJobStreamClient,
		TypeFields:                            typeFieldsClient,
		Usages:                                usagesClient,
		Variable:                              variableClient,
		Watcher:                               watcherClient,
		Webhook:                               webhookClient,
	}, nil
}
