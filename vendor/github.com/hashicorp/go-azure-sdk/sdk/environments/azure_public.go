package environments

func AzurePublic() *Environment {
	env := baseEnvironmentWithName("Public")

	env.Authorization = &Authorization{
		Audiences: []string{
			"https://management.core.windows.net",
			"https://management.azure.com",
		},
		LoginEndpoint: "https://login.microsoftonline.com",
	}
	env.ResourceManager = ResourceManagerAPI("https://management.azure.com")
	env.MicrosoftGraph = MicrosoftGraphAPI("https://graph.microsoft.com")

	env.ApiManagement = ApiManagementAPI("azure-api.net")
	env.Batch = BatchAPI("https://batch.core.windows.net")
	env.CDNFrontDoor = CDNFrontDoorAPI("azurefd.net")
	env.ContainerRegistry = ContainerRegistryAPI("azurecr.io")
	env.CosmosDB = CosmosDBAPI("documents.azure.com")
	env.DataLake = DataLakeAPI("azuredatalakestore.net")
	env.KeyVault = KeyVaultAPI("vault.azure.net")
	env.ManagedHSM = ManagedHSMAPI("https://managedhsm.azure.net", "managedhsm.azure.net")
	env.MariaDB = MariaDBAPI("mariadb.database.azure.com")
	env.MySql = MySqlAPI("mysql.database.azure.com")
	env.OperationalInsights = OperationalInsightsAPI()
	env.Postgresql = PostgresqlAPI("postgres.database.azure.com")
	env.ServiceBus = ServiceBusAPI("https://servicebus.windows.net", "servicebus.windows.net")
	env.Sql = SqlAPI("database.windows.net")
	env.Storage = StorageAPI("core.windows.net")
	env.Synapse = SynapseAPI("dev.azuresynapse.net")
	env.TrafficManager = TrafficManagerAPI("trafficmanager.net")

	return &env
}
