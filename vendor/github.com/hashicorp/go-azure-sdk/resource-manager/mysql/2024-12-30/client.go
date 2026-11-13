package v2024_12_30

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2024-12-30/advancedthreatprotectionsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2024-12-30/azureadadministrators"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2024-12-30/backupandexport"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2024-12-30/backups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2024-12-30/configurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2024-12-30/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2024-12-30/firewallrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2024-12-30/locationbasedcapabilityset"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2024-12-30/logfiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2024-12-30/longrunningbackup"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2024-12-30/longrunningbackups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2024-12-30/maintenances"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2024-12-30/openapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2024-12-30/privateendpointconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2024-12-30/privatelinkresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2024-12-30/replicas"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2024-12-30/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2024-12-30/serversmigration"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	AdvancedThreatProtectionSettings *advancedthreatprotectionsettings.AdvancedThreatProtectionSettingsClient
	AzureADAdministrators            *azureadadministrators.AzureADAdministratorsClient
	BackupAndExport                  *backupandexport.BackupAndExportClient
	Backups                          *backups.BackupsClient
	Configurations                   *configurations.ConfigurationsClient
	Databases                        *databases.DatabasesClient
	FirewallRules                    *firewallrules.FirewallRulesClient
	LocationBasedCapabilitySet       *locationbasedcapabilityset.LocationBasedCapabilitySetClient
	LogFiles                         *logfiles.LogFilesClient
	LongRunningBackup                *longrunningbackup.LongRunningBackupClient
	LongRunningBackups               *longrunningbackups.LongRunningBackupsClient
	Maintenances                     *maintenances.MaintenancesClient
	Openapis                         *openapis.OpenapisClient
	PrivateEndpointConnections       *privateendpointconnections.PrivateEndpointConnectionsClient
	PrivateLinkResources             *privatelinkresources.PrivateLinkResourcesClient
	Replicas                         *replicas.ReplicasClient
	Servers                          *servers.ServersClient
	ServersMigration                 *serversmigration.ServersMigrationClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	advancedThreatProtectionSettingsClient, err := advancedthreatprotectionsettings.NewAdvancedThreatProtectionSettingsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building AdvancedThreatProtectionSettings client: %+v", err)
	}
	configureFunc(advancedThreatProtectionSettingsClient.Client)

	azureADAdministratorsClient, err := azureadadministrators.NewAzureADAdministratorsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building AzureADAdministrators client: %+v", err)
	}
	configureFunc(azureADAdministratorsClient.Client)

	backupAndExportClient, err := backupandexport.NewBackupAndExportClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building BackupAndExport client: %+v", err)
	}
	configureFunc(backupAndExportClient.Client)

	backupsClient, err := backups.NewBackupsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Backups client: %+v", err)
	}
	configureFunc(backupsClient.Client)

	configurationsClient, err := configurations.NewConfigurationsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Configurations client: %+v", err)
	}
	configureFunc(configurationsClient.Client)

	databasesClient, err := databases.NewDatabasesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Databases client: %+v", err)
	}
	configureFunc(databasesClient.Client)

	firewallRulesClient, err := firewallrules.NewFirewallRulesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building FirewallRules client: %+v", err)
	}
	configureFunc(firewallRulesClient.Client)

	locationBasedCapabilitySetClient, err := locationbasedcapabilityset.NewLocationBasedCapabilitySetClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building LocationBasedCapabilitySet client: %+v", err)
	}
	configureFunc(locationBasedCapabilitySetClient.Client)

	logFilesClient, err := logfiles.NewLogFilesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building LogFiles client: %+v", err)
	}
	configureFunc(logFilesClient.Client)

	longRunningBackupClient, err := longrunningbackup.NewLongRunningBackupClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building LongRunningBackup client: %+v", err)
	}
	configureFunc(longRunningBackupClient.Client)

	longRunningBackupsClient, err := longrunningbackups.NewLongRunningBackupsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building LongRunningBackups client: %+v", err)
	}
	configureFunc(longRunningBackupsClient.Client)

	maintenancesClient, err := maintenances.NewMaintenancesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Maintenances client: %+v", err)
	}
	configureFunc(maintenancesClient.Client)

	openapisClient, err := openapis.NewOpenapisClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Openapis client: %+v", err)
	}
	configureFunc(openapisClient.Client)

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

	replicasClient, err := replicas.NewReplicasClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Replicas client: %+v", err)
	}
	configureFunc(replicasClient.Client)

	serversClient, err := servers.NewServersClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Servers client: %+v", err)
	}
	configureFunc(serversClient.Client)

	serversMigrationClient, err := serversmigration.NewServersMigrationClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ServersMigration client: %+v", err)
	}
	configureFunc(serversMigrationClient.Client)

	return &Client{
		AdvancedThreatProtectionSettings: advancedThreatProtectionSettingsClient,
		AzureADAdministrators:            azureADAdministratorsClient,
		BackupAndExport:                  backupAndExportClient,
		Backups:                          backupsClient,
		Configurations:                   configurationsClient,
		Databases:                        databasesClient,
		FirewallRules:                    firewallRulesClient,
		LocationBasedCapabilitySet:       locationBasedCapabilitySetClient,
		LogFiles:                         logFilesClient,
		LongRunningBackup:                longRunningBackupClient,
		LongRunningBackups:               longRunningBackupsClient,
		Maintenances:                     maintenancesClient,
		Openapis:                         openapisClient,
		PrivateEndpointConnections:       privateEndpointConnectionsClient,
		PrivateLinkResources:             privateLinkResourcesClient,
		Replicas:                         replicasClient,
		Servers:                          serversClient,
		ServersMigration:                 serversMigrationClient,
	}, nil
}
