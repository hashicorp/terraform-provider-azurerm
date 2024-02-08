package v2017_12_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2017-12-01/checknameavailability"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2017-12-01/configurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2017-12-01/configurationsupdate"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2017-12-01/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2017-12-01/firewallrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2017-12-01/locationbasedperformancetier"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2017-12-01/logfiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2017-12-01/recoverableservers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2017-12-01/replicas"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2017-12-01/serveradministrators"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2017-12-01/serverbasedperformancetier"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2017-12-01/serverrestart"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2017-12-01/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2017-12-01/serversecurityalertpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2017-12-01/virtualnetworkrules"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	CheckNameAvailability        *checknameavailability.CheckNameAvailabilityClient
	Configurations               *configurations.ConfigurationsClient
	ConfigurationsUpdate         *configurationsupdate.ConfigurationsUpdateClient
	Databases                    *databases.DatabasesClient
	FirewallRules                *firewallrules.FirewallRulesClient
	LocationBasedPerformanceTier *locationbasedperformancetier.LocationBasedPerformanceTierClient
	LogFiles                     *logfiles.LogFilesClient
	RecoverableServers           *recoverableservers.RecoverableServersClient
	Replicas                     *replicas.ReplicasClient
	ServerAdministrators         *serveradministrators.ServerAdministratorsClient
	ServerBasedPerformanceTier   *serverbasedperformancetier.ServerBasedPerformanceTierClient
	ServerRestart                *serverrestart.ServerRestartClient
	ServerSecurityAlertPolicies  *serversecurityalertpolicies.ServerSecurityAlertPoliciesClient
	Servers                      *servers.ServersClient
	VirtualNetworkRules          *virtualnetworkrules.VirtualNetworkRulesClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
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

	configurationsUpdateClient, err := configurationsupdate.NewConfigurationsUpdateClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ConfigurationsUpdate client: %+v", err)
	}
	configureFunc(configurationsUpdateClient.Client)

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

	locationBasedPerformanceTierClient, err := locationbasedperformancetier.NewLocationBasedPerformanceTierClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building LocationBasedPerformanceTier client: %+v", err)
	}
	configureFunc(locationBasedPerformanceTierClient.Client)

	logFilesClient, err := logfiles.NewLogFilesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building LogFiles client: %+v", err)
	}
	configureFunc(logFilesClient.Client)

	recoverableServersClient, err := recoverableservers.NewRecoverableServersClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building RecoverableServers client: %+v", err)
	}
	configureFunc(recoverableServersClient.Client)

	replicasClient, err := replicas.NewReplicasClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Replicas client: %+v", err)
	}
	configureFunc(replicasClient.Client)

	serverAdministratorsClient, err := serveradministrators.NewServerAdministratorsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ServerAdministrators client: %+v", err)
	}
	configureFunc(serverAdministratorsClient.Client)

	serverBasedPerformanceTierClient, err := serverbasedperformancetier.NewServerBasedPerformanceTierClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ServerBasedPerformanceTier client: %+v", err)
	}
	configureFunc(serverBasedPerformanceTierClient.Client)

	serverRestartClient, err := serverrestart.NewServerRestartClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ServerRestart client: %+v", err)
	}
	configureFunc(serverRestartClient.Client)

	serverSecurityAlertPoliciesClient, err := serversecurityalertpolicies.NewServerSecurityAlertPoliciesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ServerSecurityAlertPolicies client: %+v", err)
	}
	configureFunc(serverSecurityAlertPoliciesClient.Client)

	serversClient, err := servers.NewServersClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Servers client: %+v", err)
	}
	configureFunc(serversClient.Client)

	virtualNetworkRulesClient, err := virtualnetworkrules.NewVirtualNetworkRulesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building VirtualNetworkRules client: %+v", err)
	}
	configureFunc(virtualNetworkRulesClient.Client)

	return &Client{
		CheckNameAvailability:        checkNameAvailabilityClient,
		Configurations:               configurationsClient,
		ConfigurationsUpdate:         configurationsUpdateClient,
		Databases:                    databasesClient,
		FirewallRules:                firewallRulesClient,
		LocationBasedPerformanceTier: locationBasedPerformanceTierClient,
		LogFiles:                     logFilesClient,
		RecoverableServers:           recoverableServersClient,
		Replicas:                     replicasClient,
		ServerAdministrators:         serverAdministratorsClient,
		ServerBasedPerformanceTier:   serverBasedPerformanceTierClient,
		ServerRestart:                serverRestartClient,
		ServerSecurityAlertPolicies:  serverSecurityAlertPoliciesClient,
		Servers:                      serversClient,
		VirtualNetworkRules:          virtualNetworkRulesClient,
	}, nil
}
