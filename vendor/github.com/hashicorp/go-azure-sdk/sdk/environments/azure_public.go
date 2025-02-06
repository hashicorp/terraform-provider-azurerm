// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package environments

const AzurePublicCloud = "Public"

func AzurePublic() *Environment {
	env := baseEnvironmentWithName(AzurePublicCloud)

	env.Authorization = &Authorization{
		Audiences: []string{
			"https://management.core.windows.net",
			"https://management.azure.com",
		},
		IdentityProvider: "AAD",
		LoginEndpoint:    "https://login.microsoftonline.com",
		Tenant:           "common",
	}
	env.ResourceManager = ResourceManagerAPI("https://management.azure.com")
	env.MicrosoftGraph = MicrosoftGraphAPI("https://graph.microsoft.com")

	env.ApiManagement = ApiManagementAPI("azure-api.net")
	env.AppConfiguration = AppConfigurationAPI("azconfig.io")
	env.Attestation = AttestationAPI("https://attest.azure.net", "attest.azure.net")
	env.Batch = BatchAPI("https://batch.core.windows.net")
	env.CDNFrontDoor = CDNFrontDoorAPI("azurefd.net")
	env.ContainerRegistry = ContainerRegistryAPI("azurecr.io")
	env.CosmosDB = CosmosDBAPI("documents.azure.com").WithResourceIdentifier("https://cosmos.azure.com")
	env.DataLake = DataLakeAPI("azuredatalakestore.net").WithResourceIdentifier("https://datalake.azure.net")
	env.IoTCentral = IoTCentral("azureiotcentral.com").WithResourceIdentifier("https://apps.azureiotcentral.com")
	env.KeyVault = KeyVaultAPI("vault.azure.net").WithResourceIdentifier("https://vault.azure.net")
	env.ManagedHSM = ManagedHSMAPI("https://managedhsm.azure.net", "managedhsm.azure.net")
	env.MariaDB = MariaDBAPI("mariadb.database.azure.com").WithResourceIdentifier("https://ossrdbms-aad.database.windows.net")
	env.MySql = MySqlAPI("mysql.database.azure.com").WithResourceIdentifier("https://ossrdbms-aad.database.windows.net")
	env.OperationalInsights = OperationalInsightsAPI().WithResourceIdentifier("https://api.loganalytics.io")
	env.Postgresql = PostgresqlAPI("postgres.database.azure.com").WithResourceIdentifier("https://ossrdbms-aad.database.windows.net")
	env.ServiceBus = ServiceBusAPI("https://servicebus.windows.net", "servicebus.windows.net").WithResourceIdentifier("https://servicebus.azure.net")
	env.Sql = SqlAPI("database.windows.net").WithResourceIdentifier("https://database.windows.net")
	env.Storage = StorageAPI("core.windows.net")
	env.Synapse = SynapseAPI("dev.azuresynapse.net").WithResourceIdentifier("https://dev.azuresynapse.net")
	env.TrafficManager = TrafficManagerAPI("trafficmanager.net")

	return &env
}
