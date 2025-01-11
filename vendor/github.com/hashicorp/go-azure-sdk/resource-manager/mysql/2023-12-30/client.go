package v2023_12_30

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/advancedthreatprotectionsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/azureadadministrators"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/backupandexport"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/backups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/backupsv2"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/checknameavailability"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/configurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/firewallrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/getprivatednszonesuffix"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/locationbasedcapabilities"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/locationbasedcapability"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/logfiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/maintenances"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/operationprogress"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/serverfailover"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/servermigration"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/serverresetgtid"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/serverrestart"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/serverstart"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/serverstop"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/servervalidateestimatehighavailability"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	AdvancedThreatProtectionSettings       *advancedthreatprotectionsettings.AdvancedThreatProtectionSettingsClient
	AzureADAdministrators                  *azureadadministrators.AzureADAdministratorsClient
	BackupAndExport                        *backupandexport.BackupAndExportClient
	Backups                                *backups.BackupsClient
	BackupsV2                              *backupsv2.BackupsV2Client
	CheckNameAvailability                  *checknameavailability.CheckNameAvailabilityClient
	Configurations                         *configurations.ConfigurationsClient
	Databases                              *databases.DatabasesClient
	FirewallRules                          *firewallrules.FirewallRulesClient
	GetPrivateDnsZoneSuffix                *getprivatednszonesuffix.GetPrivateDnsZoneSuffixClient
	LocationBasedCapabilities              *locationbasedcapabilities.LocationBasedCapabilitiesClient
	LocationBasedCapability                *locationbasedcapability.LocationBasedCapabilityClient
	LogFiles                               *logfiles.LogFilesClient
	Maintenances                           *maintenances.MaintenancesClient
	OperationProgress                      *operationprogress.OperationProgressClient
	ServerFailover                         *serverfailover.ServerFailoverClient
	ServerMigration                        *servermigration.ServerMigrationClient
	ServerResetGtid                        *serverresetgtid.ServerResetGtidClient
	ServerRestart                          *serverrestart.ServerRestartClient
	ServerStart                            *serverstart.ServerStartClient
	ServerStop                             *serverstop.ServerStopClient
	ServerValidateEstimateHighAvailability *servervalidateestimatehighavailability.ServerValidateEstimateHighAvailabilityClient
	Servers                                *servers.ServersClient
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

	backupsV2Client, err := backupsv2.NewBackupsV2ClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building BackupsV2 client: %+v", err)
	}
	configureFunc(backupsV2Client.Client)

	checkNameAvailabilityClient, err := checknameavailability.NewCheckNameAvailabilityClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building CheckNameAvailability client: %+v", err)
	}
	configureFunc(checkNameAvailabilityClient.Client)

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

	getPrivateDnsZoneSuffixClient, err := getprivatednszonesuffix.NewGetPrivateDnsZoneSuffixClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building GetPrivateDnsZoneSuffix client: %+v", err)
	}
	configureFunc(getPrivateDnsZoneSuffixClient.Client)

	locationBasedCapabilitiesClient, err := locationbasedcapabilities.NewLocationBasedCapabilitiesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building LocationBasedCapabilities client: %+v", err)
	}
	configureFunc(locationBasedCapabilitiesClient.Client)

	locationBasedCapabilityClient, err := locationbasedcapability.NewLocationBasedCapabilityClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building LocationBasedCapability client: %+v", err)
	}
	configureFunc(locationBasedCapabilityClient.Client)

	logFilesClient, err := logfiles.NewLogFilesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building LogFiles client: %+v", err)
	}
	configureFunc(logFilesClient.Client)

	maintenancesClient, err := maintenances.NewMaintenancesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Maintenances client: %+v", err)
	}
	configureFunc(maintenancesClient.Client)

	operationProgressClient, err := operationprogress.NewOperationProgressClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building OperationProgress client: %+v", err)
	}
	configureFunc(operationProgressClient.Client)

	serverFailoverClient, err := serverfailover.NewServerFailoverClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ServerFailover client: %+v", err)
	}
	configureFunc(serverFailoverClient.Client)

	serverMigrationClient, err := servermigration.NewServerMigrationClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ServerMigration client: %+v", err)
	}
	configureFunc(serverMigrationClient.Client)

	serverResetGtidClient, err := serverresetgtid.NewServerResetGtidClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ServerResetGtid client: %+v", err)
	}
	configureFunc(serverResetGtidClient.Client)

	serverRestartClient, err := serverrestart.NewServerRestartClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ServerRestart client: %+v", err)
	}
	configureFunc(serverRestartClient.Client)

	serverStartClient, err := serverstart.NewServerStartClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ServerStart client: %+v", err)
	}
	configureFunc(serverStartClient.Client)

	serverStopClient, err := serverstop.NewServerStopClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ServerStop client: %+v", err)
	}
	configureFunc(serverStopClient.Client)

	serverValidateEstimateHighAvailabilityClient, err := servervalidateestimatehighavailability.NewServerValidateEstimateHighAvailabilityClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ServerValidateEstimateHighAvailability client: %+v", err)
	}
	configureFunc(serverValidateEstimateHighAvailabilityClient.Client)

	serversClient, err := servers.NewServersClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Servers client: %+v", err)
	}
	configureFunc(serversClient.Client)

	return &Client{
		AdvancedThreatProtectionSettings:       advancedThreatProtectionSettingsClient,
		AzureADAdministrators:                  azureADAdministratorsClient,
		BackupAndExport:                        backupAndExportClient,
		Backups:                                backupsClient,
		BackupsV2:                              backupsV2Client,
		CheckNameAvailability:                  checkNameAvailabilityClient,
		Configurations:                         configurationsClient,
		Databases:                              databasesClient,
		FirewallRules:                          firewallRulesClient,
		GetPrivateDnsZoneSuffix:                getPrivateDnsZoneSuffixClient,
		LocationBasedCapabilities:              locationBasedCapabilitiesClient,
		LocationBasedCapability:                locationBasedCapabilityClient,
		LogFiles:                               logFilesClient,
		Maintenances:                           maintenancesClient,
		OperationProgress:                      operationProgressClient,
		ServerFailover:                         serverFailoverClient,
		ServerMigration:                        serverMigrationClient,
		ServerResetGtid:                        serverResetGtidClient,
		ServerRestart:                          serverRestartClient,
		ServerStart:                            serverStartClient,
		ServerStop:                             serverStopClient,
		ServerValidateEstimateHighAvailability: serverValidateEstimateHighAvailabilityClient,
		Servers:                                serversClient,
	}, nil
}
