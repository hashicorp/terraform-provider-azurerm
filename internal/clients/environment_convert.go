package clients

import (
	"fmt"

	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

func toAutorestEnv(env environments.Environment) (*azure.Environment, error) {
	var resourceManagerEndpoint string
	if ptr, ok := env.ResourceManager.Endpoint(); ok {
		resourceManagerEndpoint = *ptr
	}

	var keyvaultEndpoint string
	// Key vault doesn't have endpoint defined, not sure why..
	if ptr, ok := env.KeyVault.ResourceIdentifier(); ok {
		keyvaultEndpoint = *ptr
	}

	var managedHSMEndpoint string
	if ptr, ok := env.ManagedHSM.Endpoint(); ok {
		managedHSMEndpoint = *ptr
	}

	var msGraphEndpoint string
	if ptr, ok := env.MicrosoftGraph.Endpoint(); ok {
		msGraphEndpoint = *ptr
	}

	var servicebusEndpoint string
	if ptr, ok := env.ServiceBus.Endpoint(); ok {
		servicebusEndpoint = *ptr
	}

	var batcheManagementEndpoint string
	if ptr, ok := env.Batch.Endpoint(); ok {
		batcheManagementEndpoint = *ptr
	}

	var storageEndpointSuffix string
	if ptr, ok := env.Storage.DomainSuffix(); ok {
		storageEndpointSuffix = *ptr
	}

	var cosmosDBDNSSuffix string
	if ptr, ok := env.CosmosDB.DomainSuffix(); ok {
		cosmosDBDNSSuffix = *ptr
	}

	var mariaDBDNSSuffix string
	if ptr, ok := env.MariaDB.DomainSuffix(); ok {
		mariaDBDNSSuffix = *ptr
	}

	var mysqlDBDNSSuffix string
	if ptr, ok := env.MySql.DomainSuffix(); ok {
		mysqlDBDNSSuffix = *ptr
	}

	var pgDNSSuffix string
	if ptr, ok := env.Postgresql.DomainSuffix(); ok {
		pgDNSSuffix = *ptr
	}

	var sqlDBDNSSuffix string
	if ptr, ok := env.Sql.DomainSuffix(); ok {
		sqlDBDNSSuffix = *ptr
	}

	var trafficManagerDNSSuffix string
	if ptr, ok := env.TrafficManager.DomainSuffix(); ok {
		trafficManagerDNSSuffix = *ptr
	}

	var keyvaultDNSSuffix string
	if ptr, ok := env.KeyVault.DomainSuffix(); ok {
		keyvaultDNSSuffix = *ptr
	}

	var managedHSMDNSSuffix string
	if ptr, ok := env.ManagedHSM.DomainSuffix(); ok {
		managedHSMDNSSuffix = *ptr
	}

	var servicebusSuffix string
	if ptr, ok := env.ServiceBus.DomainSuffix(); ok {
		servicebusSuffix = *ptr
	}

	var containerRegistryDNSSuffix string
	if ptr, ok := env.ContainerRegistry.DomainSuffix(); ok {
		containerRegistryDNSSuffix = *ptr
	}

	var apiManagementHostNameSuffix string
	if ptr, ok := env.ApiManagement.DomainSuffix(); ok {
		apiManagementHostNameSuffix = *ptr
	}

	var synapseEndpointSuffix string
	if ptr, ok := env.Synapse.DomainSuffix(); ok {
		synapseEndpointSuffix = *ptr
	}

	var datalakeSuffix string
	if ptr, ok := env.DataLake.DomainSuffix(); ok {
		datalakeSuffix = *ptr
	}

	var storageIdentifier string
	if ptr, ok := env.Storage.ResourceIdentifier(); ok {
		storageIdentifier = *ptr
	}

	var keyvaultIdentifier string
	if ptr, ok := env.Keyvault.ResourceIdentifier(); ok {
		keyvaultIdentifier = *ptr
	}

	var datalakeIdentifier string
	if ptr, ok := env.DataLake.ResourceIdentifier(); ok {
		datalakeIdentifier = *ptr
	}

	var batchIdentifier string
	if ptr, ok := env.Batch.ResourceIdentifier(); ok {
		batchIdentifier = *ptr
	}

	var synapseIdentifier string
	if ptr, ok := env.Synapse.ResourceIdentifier(); ok {
		synapseIdentifier = *ptr
	}

	var servicebusIdentifer string
	if ptr, ok := env.ServiceBus.ResourceIdentifier(); ok {
		servicebusIdentifer = *ptr
	}

	var operationInsightsIdentifier string
	if ptr, ok := env.OperationalInsights.ResourceIdentifier(); ok {
		operationInsightsIdentifier = *ptr
	}

	aEnv := &azure.Environment{
		Name:                         env.Name,
		ManagementPortalURL:          "", // environment only has the appId defined
		PublishSettingsURL:           "", // not available in environment
		ServiceManagementEndpoint:    "", // environment only has the appId defined
		ResourceManagerEndpoint:      resourceManagerEndpoint,
		ActiveDirectoryEndpoint:      env.Authorization.LoginEndpoint,
		GalleryEndpoint:              "", // not available in environment
		KeyVaultEndpoint:             keyvaultEndpoint,
		ManagedHSMEndpoint:           managedHSMEndpoint,
		GraphEndpoint:                "", // not available in environment
		ServiceBusEndpoint:           servicebusEndpoint,
		BatchManagementEndpoint:      batcheManagementEndpoint,
		MicrosoftGraphEndpoint:       msGraphEndpoint,
		StorageEndpointSuffix:        storageEndpointSuffix,
		CosmosDBDNSSuffix:            cosmosDBDNSSuffix,
		MariaDBDNSSuffix:             mariaDBDNSSuffix,
		MySQLDatabaseDNSSuffix:       mysqlDBDNSSuffix,
		PostgresqlDatabaseDNSSuffix:  pgDNSSuffix,
		SQLDatabaseDNSSuffix:         sqlDBDNSSuffix,
		TrafficManagerDNSSuffix:      trafficManagerDNSSuffix,
		KeyVaultDNSSuffix:            keyvaultDNSSuffix,
		ManagedHSMDNSSuffix:          managedHSMDNSSuffix,
		ServiceBusEndpointSuffix:     servicebusSuffix,
		ServiceManagementVMDNSSuffix: "", // not available in environment
		ResourceManagerVMDNSSuffix:   "", // not available in environment
		ContainerRegistryDNSSuffix:   containerRegistryDNSSuffix,
		APIManagementHostNameSuffix:  apiManagementHostNameSuffix,
		SynapseEndpointSuffix:        synapseEndpointSuffix,
		DatalakeSuffix:               datalakeSuffix,
		ResourceIdentifiers: azure.ResourceIdentifier{
			Storage:             "https://storage.azure.com/",
			Graph:               "", // not available in environment
			KeyVault:            keyvaultIdentifier,
			Datalake:            datalakeIdentifier,
			Batch:               batchIdentifier,
			Synapse:             synapseIdentifier,
			ServiceBus:          servicebusIdentifer,
			OperationalInsights: operationInsightsIdentifier,
		},
	}

	if len(env.Authorization.Audiences) > 0 {
		aEnv.TokenAudience = env.Authorization.Audiences[0]
	} else {
		return nil, fmt.Errorf("unable to find token audience for environment %q", env.Name)
	}

	return aEnv, nil
}
