// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package environments

import "github.com/hashicorp/go-azure-helpers/lang/pointer"

func applicationIdOnly(name, applicationId string) Api {
	return &ApiEndpoint{
		domainSuffix:       nil,
		endpoint:           nil,
		appId:              pointer.To(applicationId),
		name:               name,
		resourceIdentifier: nil,
	}
}

func ApiManagementAPI(domainSuffix string) Api {
	return &ApiEndpoint{
		domainSuffix:       pointer.To(domainSuffix),
		endpoint:           nil,
		appId:              pointer.To(apiManagementAppId),
		name:               "ApiManagement",
		resourceIdentifier: nil,
	}
}

func AttestationAPI(endpoint string) Api {
	return &ApiEndpoint{
		domainSuffix:       nil,
		endpoint:           pointer.To(endpoint),
		appId:              pointer.To(attestationServiceAppId),
		name:               "AttestationService",
		resourceIdentifier: nil,
	}
}

func BatchAPI(endpoint string) *ApiEndpoint {
	return &ApiEndpoint{
		domainSuffix:       nil,
		endpoint:           pointer.To(endpoint),
		appId:              pointer.To(batchAppId),
		name:               "Batch",
		resourceIdentifier: pointer.To("https://batch.core.windows.net"),
	}
}

func CDNFrontDoorAPI(domainSuffix string) Api {
	return &ApiEndpoint{
		domainSuffix:       pointer.To(domainSuffix),
		endpoint:           nil,
		appId:              pointer.To(cdnFrontDoorAppId),
		name:               "CDNFrontDoor",
		resourceIdentifier: nil,
	}
}

func ContainerRegistryAPI(domainSuffix string) Api {
	return &ApiEndpoint{
		domainSuffix:       pointer.To(domainSuffix),
		endpoint:           nil,
		appId:              pointer.To(containerRegistryAppId),
		name:               "ContainerRegistry",
		resourceIdentifier: nil,
	}
}

func CosmosDBAPI(domainSuffix string) Api {
	return &ApiEndpoint{
		domainSuffix:       pointer.To(domainSuffix),
		endpoint:           nil,
		appId:              pointer.To(cosmosDBAppId),
		name:               "AzureCosmosDb",
		resourceIdentifier: pointer.To("https://cosmos.azure.com"),
	}
}

func DataLakeAPI(domainSuffix string) Api {
	return &ApiEndpoint{
		domainSuffix:       pointer.To(domainSuffix),
		endpoint:           nil,
		appId:              pointer.To(dataLakeAppId),
		name:               "DataLake",
		resourceIdentifier: pointer.To("https://datalake.azure.net"),
	}
}

func IoTCentral(domainSuffix string) Api {
	return &ApiEndpoint{
		domainSuffix:       pointer.To(domainSuffix),
		endpoint:           nil,
		appId:              pointer.To(iotCentralAppId),
		name:               "IoTCentral",
		resourceIdentifier: pointer.To("https://apps.azureiotcentral.com"),
	}
}

func KeyVaultAPI(domainSuffix string) *ApiEndpoint {
	return &ApiEndpoint{
		domainSuffix:       pointer.To(domainSuffix),
		endpoint:           nil,
		appId:              pointer.To(keyVaultAppId),
		name:               "AzureKeyVault",
		resourceIdentifier: pointer.To("https://vault.azure.net"),
	}
}

func ManagedHSMAPI(endpoint, domainSuffix string) Api {
	return &ApiEndpoint{
		domainSuffix:       pointer.To(domainSuffix),
		endpoint:           pointer.To(endpoint),
		appId:              pointer.To(managedHSMAppId),
		name:               "ManagedHSM",
		resourceIdentifier: pointer.To("https://managedhsm.azure.net"),
	}
}

func MariaDBAPI(domainSuffix string) *ApiEndpoint {
	return &ApiEndpoint{
		domainSuffix:       pointer.To(domainSuffix),
		endpoint:           nil,
		appId:              pointer.To(ossRDBMSAppId),
		name:               "OssRdbms",
		resourceIdentifier: pointer.To("https://ossrdbms-aad.database.windows.net"),
	}
}

func MicrosoftGraphAPI(endpoint string) *ApiEndpoint {
	return &ApiEndpoint{
		domainSuffix:       nil,
		endpoint:           pointer.To(endpoint),
		appId:              pointer.To(microsoftGraphAppId),
		name:               "MicrosoftGraph",
		resourceIdentifier: pointer.To("https://graph.microsoft.com"),
	}
}

func MySqlAPI(domainSuffix string) *ApiEndpoint {
	return &ApiEndpoint{
		domainSuffix:       pointer.To(domainSuffix),
		endpoint:           nil,
		appId:              pointer.To(ossRDBMSAppId),
		name:               "OssRdbms",
		resourceIdentifier: pointer.To("https://ossrdbms-aad.database.windows.net"),
	}
}

func OperationalInsightsAPI() *ApiEndpoint {
	return &ApiEndpoint{
		domainSuffix:       nil,
		endpoint:           nil,
		appId:              pointer.To(logAnalyticsAppId),
		name:               "OperationalInsights",
		resourceIdentifier: pointer.To("https://api.loganalytics.io"),
	}
}

func PostgresqlAPI(domainSuffix string) *ApiEndpoint {
	return &ApiEndpoint{
		domainSuffix:       pointer.To(domainSuffix),
		endpoint:           nil,
		appId:              pointer.To(ossRDBMSAppId),
		name:               "OssRdbms",
		resourceIdentifier: pointer.To("https://ossrdbms-aad.database.windows.net"),
	}
}

func ResourceManagerAPI(endpoint string) *ApiEndpoint {
	return &ApiEndpoint{
		domainSuffix:       nil,
		endpoint:           pointer.To(endpoint),
		appId:              pointer.To(azureServiceManagementAppId),
		name:               "ResourceManager",
		resourceIdentifier: pointer.To("https://management.azure.com"),
	}
}

func ServiceBusAPI(endpoint, domainSuffix string) Api {
	return &ApiEndpoint{
		domainSuffix:       pointer.To(domainSuffix),
		endpoint:           pointer.To(endpoint),
		appId:              pointer.To(serviceBusAppId),
		name:               "ServiceBus",
		resourceIdentifier: pointer.To("https://servicebus.azure.net"),
	}
}

func SqlAPI(domainSuffix string) *ApiEndpoint {
	return &ApiEndpoint{
		domainSuffix:       pointer.To(domainSuffix),
		endpoint:           nil,
		appId:              pointer.To(sqlDatabaseAppId),
		name:               "AzureSqlDatabase",
		resourceIdentifier: pointer.To("https://database.windows.net"),
	}
}

func StorageAPI(domainSuffix string) *ApiEndpoint {
	return &ApiEndpoint{
		domainSuffix:       pointer.To(domainSuffix),
		endpoint:           nil,
		appId:              pointer.To(storageAppId),
		name:               "AzureStorage",
		resourceIdentifier: pointer.To("https://storage.azure.com"),
	}
}

func StorageSyncAPI(domainSuffix string) Api {
	return &ApiEndpoint{
		domainSuffix:       pointer.To(domainSuffix),
		endpoint:           nil,
		appId:              pointer.To(storageSyncAppId),
		name:               "StorageSync",
		resourceIdentifier: nil,
	}
}

func SynapseAPI(domainSuffix string) *ApiEndpoint {
	return &ApiEndpoint{
		domainSuffix:       pointer.To(domainSuffix),
		endpoint:           nil,
		appId:              pointer.To(synapseAppId),
		name:               "Synapse",
		resourceIdentifier: pointer.To("https://dev.azuresynapse.net"),
	}
}

func TrafficManagerAPI(domainSuffix string) Api {
	return &ApiEndpoint{
		domainSuffix:       pointer.To(domainSuffix),
		endpoint:           nil,
		appId:              pointer.To(trafficManagerAppId),
		name:               "TrafficManager",
		resourceIdentifier: nil,
	}
}
