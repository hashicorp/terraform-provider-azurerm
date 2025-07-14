// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resourcemanager

import "github.com/hashicorp/go-azure-sdk/sdk/client/pollers"

var longRunningOperationCustomStatuses = map[status]pollers.PollingStatus{
	// Expected/Documented Terminal Statuses
	statusCanceled:   pollers.PollingStatusCancelled,
	statusCancelled:  pollers.PollingStatusCancelled,
	statusFailed:     pollers.PollingStatusFailed,
	statusInProgress: pollers.PollingStatusInProgress,
	statusSucceeded:  pollers.PollingStatusSucceeded,

	// Unexpected Statuses below this point:

	// whilst the standard set above should be sufficient, some APIs differ from the spec and should be documented below:
	// Dashboard@2022-08-01 returns `Accepted` rather than `InProgress` during creation
	"Accepted": pollers.PollingStatusInProgress,

	// EventGrid@2022-06-15 returns `Active` rather than `InProgress` during creation
	"Active": pollers.PollingStatusInProgress,

	// NetAppVolumeReplication @ 2023-05-01 returns `AuthorizeReplication` during authorizing replication
	"AuthorizeReplication": pollers.PollingStatusInProgress,

	// VMWare @ 2022-05-01 returns `Building` rather than `InProgress` during creation
	"Building": pollers.PollingStatusInProgress,

	// NetAppVolumeReplication @ 2023-05-01 returns `BreakReplication` during breaking replication
	"BreakReplication": pollers.PollingStatusInProgress,

	// Mysql @ 2022-01-01 returns `CancelInProgress` during Update
	"CancelInProgress": pollers.PollingStatusInProgress,

	// CostManagement@2021-10-01 returns `Completed` rather than `Succeeded`: https://github.com/Azure/azure-sdk-for-go/issues/20342
	"Completed": pollers.PollingStatusSucceeded,

	// StreamAnalytics@2020-03-01 introduced `ConfiguringNetworking` as undocumented granular statuses on 2024-04-09
	"ConfiguringNetworking": pollers.PollingStatusInProgress,

	// ServiceFabricManaged @ 2021-05-01 (NodeTypes CreateOrUpdate) returns `Created` rather than `InProgress` during Creation
	"Created": pollers.PollingStatusInProgress,

	// ContainerRegistry@2019-06-01-preview returns `Creating` rather than `InProgress` during creation
	"Creating": pollers.PollingStatusInProgress,

	// StreamAnalytics@2020-03-01 introduced `CreatingVirtualMachines` as undocumented granular statuses on 2024-04-09
	"CreatingVirtualMachines": pollers.PollingStatusInProgress,

	// CosmosDB @ 2023-04-15 returns `Dequeued` rather than `InProgress` during creation/update
	"Dequeued": pollers.PollingStatusInProgress,

	// StorageSync@2020-03-01 returns `finishNewStorageSyncService` rather than `InProgress` during creation/update (https://github.com/hashicorp/go-azure-sdk/issues/565)
	"finishNewStorageSyncService": pollers.PollingStatusInProgress,

	// StorageSync@2020-03-01 returns `newManagedIdentityCredentialStep` rather than `InProgress` during creation/update (https://github.com/hashicorp/go-azure-sdk/issues/565)
	"newManagedIdentityCredentialStep": pollers.PollingStatusInProgress,

	// StorageSync@2020-03-01 returns `newPrivateDnsEntries` rather than `InProgress` during creation/update (https://github.com/hashicorp/go-azure-sdk/issues/565)
	"newPrivateDnsEntries": pollers.PollingStatusInProgress,

	// StorageSync@2020-03-01 (CloudEndpoints) returns `newReplicaGroup` rather than `InProgress` during creation/update (https://github.com/hashicorp/go-azure-sdk/issues/565)
	"newReplicaGroup": pollers.PollingStatusInProgress,

	// StorageSync@2020-03-01 returns `notifySyncServicePartition` rather than `InProgress` during creation
	// polling after StorageSyncServicesCreate: `result.Status` was nil/empty - `op.Status` was "notifySyncServicePartition" / `op.Properties.ProvisioningState` was ""
	"notifySyncServicePartition": pollers.PollingStatusInProgress,

	// NetApp @ 2023-05-01 (Volume Update) returns `Patching` during Update
	"Patching": pollers.PollingStatusInProgress,

	// AnalysisServices @ 2017-08-01 (Servers Suspend) returns `Pausing` during update
	"Pausing": pollers.PollingStatusInProgress,

	// ContainerInstance @ 2023-05-01 returns `Pending` during creation/update
	"Pending": pollers.PollingStatusInProgress,

	// SAPVirtualInstance @ 2023-04-01 returns `Preparing System Configuration` during Creation
	"Preparing System Configuration": pollers.PollingStatusInProgress,

	// AnalysisServices @ 2017-08-01 (Servers) returns `Provisioning` during Creation
	"Provisioning": pollers.PollingStatusInProgress,

	// Resources @ 2020-10-01 (DeploymentScripts) returns `ProvisioningResources` during Creation
	"ProvisioningResources": pollers.PollingStatusInProgress,

	// AnalysisServices @ 2017-08-01 (Servers Resume) returns `Resuming` during Update
	"Resuming": pollers.PollingStatusInProgress,

	// HealthcareApis @ 2022-12-01 returns `Requested` during Creation
	"Requested": pollers.PollingStatusInProgress,

	// SignalR@2022-02-01 returns `Running` rather than `InProgress` during creation
	"Running": pollers.PollingStatusInProgress,

	// AnalysisServices @ 2017-08-01 (Servers Suspend) returns `Scaling` during Update
	"Scaling": pollers.PollingStatusInProgress,

	// StreamAnalytics@2020-03-01 introduced `SettingUpStreamingRuntime` as undocumented granular statuses on 2024-04-09
	"SettingUpStreamingRuntime": pollers.PollingStatusInProgress,

	// KubernetesConfiguration@2022-11-01 returns `Updating` rather than `InProgress` during update
	"Updating": pollers.PollingStatusInProgress,

	// HealthBot @ 2022-08-08 (HealthBots CreateOrUpdate) returns `Working` during Creation
	"Working": pollers.PollingStatusInProgress,

	// StorageSync@2020-03-01 returns `validateInput` rather than `InProgress` during creation/update (https://github.com/hashicorp/go-azure-sdk/issues/565)
	"validateInput": pollers.PollingStatusInProgress,

	// EventGrid @ 2022-06-15 returns `AwaitingManualAction` while waiting for manual validation of a webhook (https://github.com/hashicorp/terraform-provider-azurerm/issues/25689)
	"AwaitingManualAction": pollers.PollingStatusInProgress,
}
