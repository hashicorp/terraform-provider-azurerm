package environments

import (
	"strings"

	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/manicminer/hamilton/environments"
)

// EnvironmentFromAzureEnvironment converts an Autorest azure.Environment into a Hamilton environments.Environment struct
func EnvironmentFromAzureEnvironment(azureEnv azure.Environment) environments.Environment {
	return environments.Environment{
		AzureADEndpoint: environments.AzureADEndpoint(strings.TrimSuffix(azureEnv.ActiveDirectoryEndpoint, "/")),

		AadGraph: environments.Api{
			AppId:    environments.PublishedApis["AzureActiveDirectoryGraph"],
			Endpoint: environments.ApiEndpoint(strings.TrimSuffix(azureEnv.GraphEndpoint, "/")),
		},

		BatchManagement: environments.Api{
			AppId:    environments.PublishedApis["AzureBatch"],
			Endpoint: environments.ApiEndpoint(strings.TrimSuffix(azureEnv.BatchManagementEndpoint, "/")),
		},

		DataLake: environments.Api{
			AppId:    environments.PublishedApis["AzureDataLake"],
			Endpoint: environments.ApiEndpoint(strings.TrimSuffix(azureEnv.ResourceIdentifiers.Datalake, "/")),
		},

		KeyVault: environments.Api{
			AppId:    environments.PublishedApis["AzureKeyVault"],
			Endpoint: environments.ApiEndpoint(strings.TrimSuffix(azureEnv.KeyVaultEndpoint, "/")),
		},

		OperationalInsights: environments.Api{
			AppId:    environments.PublishedApis["LogAnalytics"],
			Endpoint: environments.ApiEndpoint(strings.TrimSuffix(azureEnv.ResourceIdentifiers.OperationalInsights, "/")),
		},

		OSSRDBMS: environments.Api{
			AppId:    environments.PublishedApis["OssRdbms"],
			Endpoint: environments.ApiEndpoint(strings.TrimSuffix(azureEnv.ResourceIdentifiers.OSSRDBMS, "/")),
		},

		ResourceManager: environments.Api{
			AppId:    environments.PublishedApis["AzureServiceManagement"],
			Endpoint: environments.ApiEndpoint(strings.TrimSuffix(azureEnv.ResourceManagerEndpoint, "/")),
		},

		ServiceBus: environments.Api{
			AppId:    environments.PublishedApis["AzureServiceBus"],
			Endpoint: environments.ApiEndpoint(strings.TrimSuffix(azureEnv.ServiceBusEndpoint, "/")),
		},

		ServiceManagement: environments.Api{
			AppId:    environments.PublishedApis["AzureServiceManagement"],
			Endpoint: environments.ApiEndpoint(strings.TrimSuffix(azureEnv.ServiceManagementEndpoint, "/")),
		},

		SQLDatabase: environments.Api{
			AppId:    environments.PublishedApis["AzureSqlDatabase"],
			Endpoint: environments.ApiEndpoint(strings.TrimSuffix(azureEnv.ResourceIdentifiers.SQLDatabase, "/")),
		},

		Storage: environments.Api{
			AppId:    environments.PublishedApis["AzureStorage"],
			Endpoint: environments.ApiEndpoint(strings.TrimSuffix(azureEnv.ResourceIdentifiers.Storage, "/")),
		},

		Synapse: environments.Api{
			AppId:    environments.PublishedApis["AzureSynapseGateway"],
			Endpoint: environments.ApiEndpoint(strings.TrimSuffix(azureEnv.ResourceIdentifiers.Synapse, "/")),
		},
	}
}
