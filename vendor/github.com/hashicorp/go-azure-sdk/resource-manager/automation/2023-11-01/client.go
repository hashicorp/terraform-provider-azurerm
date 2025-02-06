package v2023_11_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/activity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/automationaccount"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/certificate"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/connection"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/connectiontype"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/credential"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/dscconfiguration"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/dscnodeconfiguration"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/hybridrunbookworker"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/hybridrunbookworkergroup"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/job"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/jobschedule"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/jobstream"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/linkedworkspace"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/listkeys"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/module"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/objectdatatypes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/operations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/python2package"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/python3package"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/runbook"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/runbookdraft"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/schedule"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/softwareupdateconfigurationmachinerun"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/softwareupdateconfigurationrun"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/sourcecontrol"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/sourcecontrolsyncjob"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/sourcecontrolsyncjobstreams"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/statistics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/testjob"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/testjobstream"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/typefields"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/usages"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/variable"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	Activity                              *activity.ActivityClient
	AutomationAccount                     *automationaccount.AutomationAccountClient
	Certificate                           *certificate.CertificateClient
	Connection                            *connection.ConnectionClient
	ConnectionType                        *connectiontype.ConnectionTypeClient
	Credential                            *credential.CredentialClient
	DscConfiguration                      *dscconfiguration.DscConfigurationClient
	DscNodeConfiguration                  *dscnodeconfiguration.DscNodeConfigurationClient
	HybridRunbookWorker                   *hybridrunbookworker.HybridRunbookWorkerClient
	HybridRunbookWorkerGroup              *hybridrunbookworkergroup.HybridRunbookWorkerGroupClient
	Job                                   *job.JobClient
	JobSchedule                           *jobschedule.JobScheduleClient
	JobStream                             *jobstream.JobStreamClient
	LinkedWorkspace                       *linkedworkspace.LinkedWorkspaceClient
	ListKeys                              *listkeys.ListKeysClient
	Module                                *module.ModuleClient
	ObjectDataTypes                       *objectdatatypes.ObjectDataTypesClient
	Operations                            *operations.OperationsClient
	Python2Package                        *python2package.Python2PackageClient
	Python3Package                        *python3package.Python3PackageClient
	Runbook                               *runbook.RunbookClient
	RunbookDraft                          *runbookdraft.RunbookDraftClient
	Schedule                              *schedule.ScheduleClient
	SoftwareUpdateConfigurationMachineRun *softwareupdateconfigurationmachinerun.SoftwareUpdateConfigurationMachineRunClient
	SoftwareUpdateConfigurationRun        *softwareupdateconfigurationrun.SoftwareUpdateConfigurationRunClient
	SourceControl                         *sourcecontrol.SourceControlClient
	SourceControlSyncJob                  *sourcecontrolsyncjob.SourceControlSyncJobClient
	SourceControlSyncJobStreams           *sourcecontrolsyncjobstreams.SourceControlSyncJobStreamsClient
	Statistics                            *statistics.StatisticsClient
	TestJob                               *testjob.TestJobClient
	TestJobStream                         *testjobstream.TestJobStreamClient
	TypeFields                            *typefields.TypeFieldsClient
	Usages                                *usages.UsagesClient
	Variable                              *variable.VariableClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	activityClient, err := activity.NewActivityClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Activity client: %+v", err)
	}
	configureFunc(activityClient.Client)

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

	dscConfigurationClient, err := dscconfiguration.NewDscConfigurationClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building DscConfiguration client: %+v", err)
	}
	configureFunc(dscConfigurationClient.Client)

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

	scheduleClient, err := schedule.NewScheduleClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Schedule client: %+v", err)
	}
	configureFunc(scheduleClient.Client)

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

	return &Client{
		Activity:                              activityClient,
		AutomationAccount:                     automationAccountClient,
		Certificate:                           certificateClient,
		Connection:                            connectionClient,
		ConnectionType:                        connectionTypeClient,
		Credential:                            credentialClient,
		DscConfiguration:                      dscConfigurationClient,
		DscNodeConfiguration:                  dscNodeConfigurationClient,
		HybridRunbookWorker:                   hybridRunbookWorkerClient,
		HybridRunbookWorkerGroup:              hybridRunbookWorkerGroupClient,
		Job:                                   jobClient,
		JobSchedule:                           jobScheduleClient,
		JobStream:                             jobStreamClient,
		LinkedWorkspace:                       linkedWorkspaceClient,
		ListKeys:                              listKeysClient,
		Module:                                moduleClient,
		ObjectDataTypes:                       objectDataTypesClient,
		Operations:                            operationsClient,
		Python2Package:                        python2PackageClient,
		Python3Package:                        python3PackageClient,
		Runbook:                               runbookClient,
		RunbookDraft:                          runbookDraftClient,
		Schedule:                              scheduleClient,
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
	}, nil
}
