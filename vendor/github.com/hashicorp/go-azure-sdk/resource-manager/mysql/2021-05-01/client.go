package v2021_05_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2021-05-01/backups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2021-05-01/checknameavailability"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2021-05-01/configurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2021-05-01/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2021-05-01/firewallrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2021-05-01/getprivatednszonesuffix"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2021-05-01/locationbasedcapabilities"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2021-05-01/serverfailover"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2021-05-01/serverrestart"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2021-05-01/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2021-05-01/serverstart"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2021-05-01/serverstop"
)

type Client struct {
	Backups                   *backups.BackupsClient
	CheckNameAvailability     *checknameavailability.CheckNameAvailabilityClient
	Configurations            *configurations.ConfigurationsClient
	Databases                 *databases.DatabasesClient
	FirewallRules             *firewallrules.FirewallRulesClient
	GetPrivateDnsZoneSuffix   *getprivatednszonesuffix.GetPrivateDnsZoneSuffixClient
	LocationBasedCapabilities *locationbasedcapabilities.LocationBasedCapabilitiesClient
	ServerFailover            *serverfailover.ServerFailoverClient
	ServerRestart             *serverrestart.ServerRestartClient
	ServerStart               *serverstart.ServerStartClient
	ServerStop                *serverstop.ServerStopClient
	Servers                   *servers.ServersClient
}

func NewClientWithBaseURI(endpoint string, configureAuthFunc func(c *autorest.Client)) Client {

	backupsClient := backups.NewBackupsClientWithBaseURI(endpoint)
	configureAuthFunc(&backupsClient.Client)

	checkNameAvailabilityClient := checknameavailability.NewCheckNameAvailabilityClientWithBaseURI(endpoint)
	configureAuthFunc(&checkNameAvailabilityClient.Client)

	configurationsClient := configurations.NewConfigurationsClientWithBaseURI(endpoint)
	configureAuthFunc(&configurationsClient.Client)

	databasesClient := databases.NewDatabasesClientWithBaseURI(endpoint)
	configureAuthFunc(&databasesClient.Client)

	firewallRulesClient := firewallrules.NewFirewallRulesClientWithBaseURI(endpoint)
	configureAuthFunc(&firewallRulesClient.Client)

	getPrivateDnsZoneSuffixClient := getprivatednszonesuffix.NewGetPrivateDnsZoneSuffixClientWithBaseURI(endpoint)
	configureAuthFunc(&getPrivateDnsZoneSuffixClient.Client)

	locationBasedCapabilitiesClient := locationbasedcapabilities.NewLocationBasedCapabilitiesClientWithBaseURI(endpoint)
	configureAuthFunc(&locationBasedCapabilitiesClient.Client)

	serverFailoverClient := serverfailover.NewServerFailoverClientWithBaseURI(endpoint)
	configureAuthFunc(&serverFailoverClient.Client)

	serverRestartClient := serverrestart.NewServerRestartClientWithBaseURI(endpoint)
	configureAuthFunc(&serverRestartClient.Client)

	serverStartClient := serverstart.NewServerStartClientWithBaseURI(endpoint)
	configureAuthFunc(&serverStartClient.Client)

	serverStopClient := serverstop.NewServerStopClientWithBaseURI(endpoint)
	configureAuthFunc(&serverStopClient.Client)

	serversClient := servers.NewServersClientWithBaseURI(endpoint)
	configureAuthFunc(&serversClient.Client)

	return Client{
		Backups:                   &backupsClient,
		CheckNameAvailability:     &checkNameAvailabilityClient,
		Configurations:            &configurationsClient,
		Databases:                 &databasesClient,
		FirewallRules:             &firewallRulesClient,
		GetPrivateDnsZoneSuffix:   &getPrivateDnsZoneSuffixClient,
		LocationBasedCapabilities: &locationBasedCapabilitiesClient,
		ServerFailover:            &serverFailoverClient,
		ServerRestart:             &serverRestartClient,
		ServerStart:               &serverStartClient,
		ServerStop:                &serverStopClient,
		Servers:                   &serversClient,
	}
}
