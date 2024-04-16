// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package environments

const AzureChinaCloud = "China"

func AzureChina() *Environment {
	env := baseEnvironmentWithName(AzureChinaCloud)

	env.Authorization = &Authorization{
		Audiences: []string{
			"https://management.core.chinacloudapi.cn",
			"https://management.chinacloudapi.cn",
		},
		IdentityProvider: "AAD",
		LoginEndpoint:    "https://login.chinacloudapi.cn",
		Tenant:           "common",
	}
	env.ResourceManager = ResourceManagerAPI("https://management.chinacloudapi.cn")
	env.MicrosoftGraph = MicrosoftGraphAPI("https://microsoftgraph.chinacloudapi.cn")

	env.ApiManagement = ApiManagementAPI("azure-api.cn")
	env.AppConfiguration = AppConfigurationAPI("azconfig.azure.cn")
	env.Batch = BatchAPI("https://batch.chinacloudapi.cn")
	env.ContainerRegistry = ContainerRegistryAPI("azurecr.cn")
	env.CosmosDB = CosmosDBAPI("documents.azure.cn")
	env.KeyVault = KeyVaultAPI("vault.azure.cn").WithResourceIdentifier("https://vault.azure.cn")
	env.ManagedHSM = ManagedHSMAPI("https://managedhsm.azure.cn", "managedhsm.azure.cn")
	env.MariaDB = MariaDBAPI("mariadb.database.chinacloudapi.cn").WithResourceIdentifier("https://ossrdbms-aad.database.chinacloudapi.cn")
	env.MySql = MySqlAPI("mysql.database.chinacloudapi.cn").WithResourceIdentifier("https://ossrdbms-aad.database.chinacloudapi.cn")
	env.OperationalInsights = OperationalInsightsAPI().WithResourceIdentifier("https://api.loganalytics.azure.cn")
	env.Postgresql = PostgresqlAPI("postgres.database.chinacloudapi.cn").WithResourceIdentifier("https://ossrdbms-aad.database.chinacloudapi.cn")
	env.ServiceBus = ServiceBusAPI("https://servicebus.chinacloudapi.cn", "servicebus.chinacloudapi.cn").WithResourceIdentifier("https://servicebus.chinacloudapi.cn")
	env.Sql = SqlAPI("database.chinacloudapi.cn").WithResourceIdentifier("https://database.chinacloudapi.cn")
	env.Storage = StorageAPI("core.chinacloudapi.cn").WithResourceIdentifier("https://storage.azure.com")
	env.Synapse = SynapseAPI("dev.azuresynapse.azure.cn").WithResourceIdentifier("https://dev.azuresynapse.azure.cn")
	env.TrafficManager = TrafficManagerAPI("trafficmanager.cn")

	// Services not currently available: Attestation, CDNFrontDoor, DataLake, IOTCentral

	return &env
}
