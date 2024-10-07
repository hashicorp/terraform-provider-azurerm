// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package environments

const AzureUSGovernmentCloud = "USGovernment"

func AzureUSGovernment() *Environment {
	env := baseEnvironmentWithName(AzureUSGovernmentCloud)

	env.Authorization = &Authorization{
		Audiences: []string{
			"https://management.core.usgovcloudapi.net",
			"https://management.usgovcloudapi.net",
		},
		IdentityProvider: "AAD",
		LoginEndpoint:    "https://login.microsoftonline.us",
		Tenant:           "common",
	}
	env.ResourceManager = ResourceManagerAPI("https://management.usgovcloudapi.net")
	env.MicrosoftGraph = MicrosoftGraphAPI("https://graph.microsoft.us")

	env.ApiManagement = ApiManagementAPI("azure-api.us")
	env.AppConfiguration = AppConfigurationAPI("azconfig.azure.us")
	env.Batch = BatchAPI("https://batch.core.usgovcloudapi.net")
	env.ContainerRegistry = ContainerRegistryAPI("azurecr.us")
	env.CosmosDB = CosmosDBAPI("documents.azure.us")
	env.KeyVault = KeyVaultAPI("vault.usgovcloudapi.net").WithResourceIdentifier("https://vault.usgovcloudapi.net")
	env.ManagedHSM = ManagedHSMAPI("https://managedhsm.usgovcloudapi.net", "managedhsm.usgovcloudapi.net")
	env.MariaDB = MariaDBAPI("mariadb.database.usgovcloudapi.net").WithResourceIdentifier("https://ossrdbms-aad.database.usgovcloudapi.net")
	env.MySql = MySqlAPI("mysql.database.usgovcloudapi.net").WithResourceIdentifier("https://ossrdbms-aad.database.usgovcloudapi.net")
	env.OperationalInsights = OperationalInsightsAPI().WithResourceIdentifier("https://api.loganalytics.us")
	env.Postgresql = PostgresqlAPI("postgres.database.usgovcloudapi.net").WithResourceIdentifier("https://ossrdbms-aad.database.usgovcloudapi.net")
	env.ServiceBus = ServiceBusAPI("https://servicebus.usgovcloudapi.net", "servicebus.usgovcloudapi.net").WithResourceIdentifier("https://servicebus.usgovcloudapi.net")
	env.Sql = SqlAPI("database.usgovcloudapi.net").WithResourceIdentifier("https://database.usgovcloudapi.net")
	env.Storage = StorageAPI("core.usgovcloudapi.net")
	env.Synapse = SynapseAPI("dev.azuresynapse.usgovcloudapi.net").WithResourceIdentifier("https://dev.azuresynapse.usgovcloudapi.net")
	env.TrafficManager = TrafficManagerAPI("usgovtrafficmanager.net")

	// Services not currently available: Attestation, CDNFrontDoor, DataLake, IOTCentral

	return &env
}

func AzureUSGovernmentL5() *Environment {
	// L5 is Azure Government with a different Microsoft Graph endpoint
	env := AzureUSGovernment()
	env.Name = "USGovernmentL5"
	env.MicrosoftGraph = MicrosoftGraphAPI("https://dod-graph.microsoft.us")
	return env
}
