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
	env.ResourceManager = ResourceManagerAPI("https://management.chinacloudapi.cn").withResourceIdentifier("https://management.chinacloudapi.cn")
	env.MicrosoftGraph = MicrosoftGraphAPI("https://microsoftgraph.chinacloudapi.cn").withResourceIdentifier("https://microsoftgraph.chinacloudapi.cn")

	// DataLake, ManagedHSM and OperationalInsights are not available
	env.ApiManagement = ApiManagementAPI("azure-api.cn")
	env.Batch = BatchAPI("https://batch.chinacloudapi.cn").withResourceIdentifier("https://batch.chinacloudapi.cn")
	env.ContainerRegistry = ContainerRegistryAPI("azurecr.cn")
	env.CosmosDB = CosmosDBAPI("documents.azure.cn")
	env.KeyVault = KeyVaultAPI("vault.azure.cn").withResourceIdentifier("https://vault.azure.cn")
	env.MariaDB = MariaDBAPI("mariadb.database.chinacloudapi.cn").withResourceIdentifier("https://ossrdbms-aad.database.chinacloudapi.cn")
	env.MySql = MySqlAPI("mysql.database.chinacloudapi.cn").withResourceIdentifier("https://ossrdbms-aad.database.chinacloudapi.cn")
	env.OperationalInsights = OperationalInsightsAPI().withResourceIdentifier("https://api.loganalytics.azure.cn")
	env.Postgresql = PostgresqlAPI("postgres.database.chinacloudapi.cn").withResourceIdentifier("https://ossrdbms-aad.database.chinacloudapi.cn")
	env.ServiceBus = ServiceBusAPI("https://servicebus.chinacloudapi.cn", "servicebus.chinacloudapi.cn")
	env.Sql = SqlAPI("database.chinacloudapi.cn").withResourceIdentifier("https://database.chinacloudapi.cn")
	env.Storage = StorageAPI("core.chinacloudapi.cn").withResourceIdentifier("https://core.chinacloudapi.cn")
	env.Synapse = SynapseAPI("dev.azuresynapse.azure.cn")
	env.TrafficManager = TrafficManagerAPI("trafficmanager.cn")

	// @tombuildsstuff: DataLake doesn't appear to be available?

	// Managed HSM expected "H2 2023" per:
	// https://azure.microsoft.com/en-gb/explore/global-infrastructure/products-by-region/?regions=china-non-regional,china-east,china-east-2,china-east-3,china-north,china-north-2,china-north-3&products=all
	// presumably this'll be
	// env.ManagedHSM = ManagedHSMAPI("https://managedhsm.azure.cn", "managedhsm.azure.cn")

	return &env
}
