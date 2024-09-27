// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package environments

import "github.com/hashicorp/go-azure-helpers/lang/pointer"

func applicationIdOnly(name, applicationId string) *ApiEndpoint {
	return &ApiEndpoint{
		domainSuffix:       nil,
		endpoint:           nil,
		appId:              pointer.To(applicationId),
		name:               name,
		resourceIdentifier: nil,
	}
}

func ApiManagementAPI(domainSuffix string) *ApiEndpoint {
	return &ApiEndpoint{
		domainSuffix:       pointer.To(domainSuffix),
		endpoint:           nil,
		appId:              pointer.To(apiManagementAppId),
		name:               "ApiManagement",
		resourceIdentifier: nil,
	}
}

func AppConfigurationAPI(domainSuffix string) *ApiEndpoint {
	return &ApiEndpoint{
		domainSuffix:       pointer.To(domainSuffix),
		endpoint:           nil,
		appId:              pointer.To(appConfigurationAppId),
		name:               "AppConfiguration",
		resourceIdentifier: nil,
	}
}

func AttestationAPI(endpoint, domainSuffix string) *ApiEndpoint {
	// endpoint and resource ID are the same, only the resource ID is returned in metadata
	return &ApiEndpoint{
		domainSuffix:       pointer.To(domainSuffix),
		endpoint:           pointer.To(endpoint),
		appId:              pointer.To(attestationServiceAppId),
		name:               "AttestationService",
		resourceIdentifier: pointer.To(endpoint),
	}
}

func BatchAPI(resourceId string) *ApiEndpoint {
	// endpoint and resource ID are the same, only the resource ID is returned in metadata
	return &ApiEndpoint{
		domainSuffix:       nil,
		endpoint:           pointer.To(resourceId),
		appId:              pointer.To(batchAppId),
		name:               "Batch",
		resourceIdentifier: pointer.To(resourceId),
	}
}

func CDNFrontDoorAPI(domainSuffix string) *ApiEndpoint {
	return &ApiEndpoint{
		domainSuffix:       pointer.To(domainSuffix),
		endpoint:           nil,
		appId:              pointer.To(cdnFrontDoorAppId),
		name:               "CDNFrontDoor",
		resourceIdentifier: nil,
	}
}

func ContainerRegistryAPI(domainSuffix string) *ApiEndpoint {
	return &ApiEndpoint{
		domainSuffix:       pointer.To(domainSuffix),
		endpoint:           nil,
		appId:              pointer.To(containerRegistryAppId),
		name:               "ContainerRegistry",
		resourceIdentifier: nil,
	}
}

func CosmosDBAPI(domainSuffix string) *ApiEndpoint {
	return &ApiEndpoint{
		domainSuffix:       pointer.To(domainSuffix),
		endpoint:           nil,
		appId:              pointer.To(cosmosDBAppId),
		name:               "AzureCosmosDb",
		resourceIdentifier: nil,
	}
}

func DataLakeAPI(domainSuffix string) *ApiEndpoint {
	return &ApiEndpoint{
		domainSuffix:       pointer.To(domainSuffix),
		endpoint:           nil,
		appId:              pointer.To(dataLakeAppId),
		name:               "DataLake",
		resourceIdentifier: nil,
	}
}

func IoTCentral(domainSuffix string) *ApiEndpoint {
	return &ApiEndpoint{
		domainSuffix:       pointer.To(domainSuffix),
		endpoint:           nil,
		appId:              pointer.To(iotCentralAppId),
		name:               "IoTCentral",
		resourceIdentifier: nil,
	}
}

func KeyVaultAPI(domainSuffix string) *ApiEndpoint {
	return &ApiEndpoint{
		domainSuffix:       pointer.To(domainSuffix),
		endpoint:           nil,
		appId:              pointer.To(keyVaultAppId),
		name:               "AzureKeyVault",
		resourceIdentifier: nil,
	}
}

func ManagedHSMAPI(endpoint, domainSuffix string) *ApiEndpoint {
	// endpoint and resource ID are the same, only the domainSuffix is returned in metadata
	return &ApiEndpoint{
		domainSuffix:       pointer.To(domainSuffix),
		endpoint:           pointer.To(endpoint),
		appId:              pointer.To(managedHSMAppId),
		name:               "ManagedHSM",
		resourceIdentifier: pointer.To(endpoint),
	}
}

func MariaDBAPI(domainSuffix string) *ApiEndpoint {
	return &ApiEndpoint{
		domainSuffix:       pointer.To(domainSuffix),
		endpoint:           nil,
		appId:              pointer.To(ossRDBMSAppId),
		name:               "OssRdbms",
		resourceIdentifier: nil,
	}
}

func MicrosoftGraphAPI(resourceId string) *ApiEndpoint {
	// endpoint and resource ID are the same, only the resource ID is returned in metadata
	return &ApiEndpoint{
		endpoint:           pointer.To(resourceId),
		appId:              pointer.To(microsoftGraphAppId),
		name:               "MicrosoftGraph",
		resourceIdentifier: pointer.To(resourceId),
	}
}

func MySqlAPI(domainSuffix string) *ApiEndpoint {
	return &ApiEndpoint{
		domainSuffix:       pointer.To(domainSuffix),
		endpoint:           nil,
		appId:              pointer.To(ossRDBMSAppId),
		name:               "OssRdbms",
		resourceIdentifier: nil,
	}
}

func OperationalInsightsAPI() *ApiEndpoint {
	return &ApiEndpoint{
		domainSuffix:       nil,
		endpoint:           nil,
		appId:              pointer.To(logAnalyticsAppId),
		name:               "OperationalInsights",
		resourceIdentifier: nil,
	}
}

func PostgresqlAPI(domainSuffix string) *ApiEndpoint {
	return &ApiEndpoint{
		domainSuffix:       pointer.To(domainSuffix),
		endpoint:           nil,
		appId:              pointer.To(ossRDBMSAppId),
		name:               "OssRdbms",
		resourceIdentifier: nil,
	}
}

func ResourceManagerAPI(endpoint string) *ApiEndpoint {
	// endpoint and resource ID are the same, only the endpoint is returned in metadata
	return &ApiEndpoint{
		domainSuffix:       nil,
		endpoint:           pointer.To(endpoint),
		appId:              pointer.To(azureServiceManagementAppId),
		name:               "ResourceManager",
		resourceIdentifier: pointer.To(endpoint),
	}
}

func ServiceBusAPI(endpoint, domainSuffix string) Api {
	return &ApiEndpoint{
		domainSuffix:       pointer.To(domainSuffix),
		endpoint:           pointer.To(endpoint),
		appId:              pointer.To(serviceBusAppId),
		name:               "ServiceBus",
		resourceIdentifier: nil,
	}
}

func SqlAPI(domainSuffix string) *ApiEndpoint {
	return &ApiEndpoint{
		domainSuffix:       pointer.To(domainSuffix),
		endpoint:           nil,
		appId:              pointer.To(sqlDatabaseAppId),
		name:               "AzureSqlDatabase",
		resourceIdentifier: nil,
	}
}

func StorageAPI(domainSuffix string) *ApiEndpoint {
	// The default resource identifier for Azure Storage is the same for all public and sovereign clouds. This can be
	// changed to scope the token for authorizing against a single storage account.
	// https://learn.microsoft.com/en-us/azure/storage/blobs/authorize-access-azure-active-directory#microsoft-authentication-library-msal
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
		resourceIdentifier: nil,
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
