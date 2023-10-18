package v2022_01_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2022-01-01/azureadadministrators"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2022-01-01/backups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2022-01-01/checknameavailability"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2022-01-01/configurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2022-01-01/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2022-01-01/firewallrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2022-01-01/getprivatednszonesuffix"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2022-01-01/locationbasedcapabilities"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2022-01-01/logfiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2022-01-01/serverfailover"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2022-01-01/serverrestart"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2022-01-01/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2022-01-01/serverstart"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2022-01-01/serverstop"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	AzureADAdministrators     *azureadadministrators.AzureADAdministratorsClient
	Backups                   *backups.BackupsClient
	CheckNameAvailability     *checknameavailability.CheckNameAvailabilityClient
	Configurations            *configurations.ConfigurationsClient
	Databases                 *databases.DatabasesClient
	FirewallRules             *firewallrules.FirewallRulesClient
	GetPrivateDnsZoneSuffix   *getprivatednszonesuffix.GetPrivateDnsZoneSuffixClient
	LocationBasedCapabilities *locationbasedcapabilities.LocationBasedCapabilitiesClient
	LogFiles                  *logfiles.LogFilesClient
	ServerFailover            *serverfailover.ServerFailoverClient
	ServerRestart             *serverrestart.ServerRestartClient
	ServerStart               *serverstart.ServerStartClient
	ServerStop                *serverstop.ServerStopClient
	Servers                   *servers.ServersClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	azureADAdministratorsClient, err := azureadadministrators.NewAzureADAdministratorsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building AzureADAdministrators client: %+v", err)
	}
	configureFunc(azureADAdministratorsClient.Client)

	backupsClient, err := backups.NewBackupsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Backups client: %+v", err)
	}
	configureFunc(backupsClient.Client)

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

	logFilesClient, err := logfiles.NewLogFilesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building LogFiles client: %+v", err)
	}
	configureFunc(logFilesClient.Client)

	serverFailoverClient, err := serverfailover.NewServerFailoverClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ServerFailover client: %+v", err)
	}
	configureFunc(serverFailoverClient.Client)

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

	serversClient, err := servers.NewServersClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Servers client: %+v", err)
	}
	configureFunc(serversClient.Client)

	return &Client{
		AzureADAdministrators:     azureADAdministratorsClient,
		Backups:                   backupsClient,
		CheckNameAvailability:     checkNameAvailabilityClient,
		Configurations:            configurationsClient,
		Databases:                 databasesClient,
		FirewallRules:             firewallRulesClient,
		GetPrivateDnsZoneSuffix:   getPrivateDnsZoneSuffixClient,
		LocationBasedCapabilities: locationBasedCapabilitiesClient,
		LogFiles:                  logFilesClient,
		ServerFailover:            serverFailoverClient,
		ServerRestart:             serverRestartClient,
		ServerStart:               serverStartClient,
		ServerStop:                serverStopClient,
		Servers:                   serversClient,
	}, nil
}
