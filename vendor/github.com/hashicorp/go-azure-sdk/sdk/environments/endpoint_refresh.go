// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package environments

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/internal/metadata"
)

func (e *Environment) RefreshMetaDataFromEndpoint(ctx context.Context) error {
	endpoint, ok := e.ResourceManager.Endpoint()
	if !ok {
		return fmt.Errorf("refreshing MetaData from Endpoint: no `ResourceManager` endpoint was defined")
	}

	client := metadata.NewClientWithEndpoint(*endpoint)
	config, err := client.GetMetaData(ctx)
	if err != nil {
		return fmt.Errorf("retrieving MetaData from endpoint: %+v", err)
	}

	if err := e.updateFromMetaData(config); err != nil {
		return fmt.Errorf("updating Environment from MetaData: %+v", err)
	}

	return nil
}

func (e *Environment) updateFromMetaData(config *metadata.MetaData) error {
	// The following supported services are missing from metadata and cannot be configured:
	// - API Management (domain suffix is missing)
	// - App Configuration (domain suffix and resource identifier are missing)
	// - CosmosDB (domain suffix is missing)
	// - IOT Central (domain suffix and resource identifier are missing)
	// - Service Bus (domain suffix and resource identifier are missing)
	// - Traffic Manager (domain suffix is missing)

	if e.Authorization == nil {
		e.Authorization = &Authorization{}
	}
	if len(config.Authentication.Audiences) > 0 {
		e.Authorization.Audiences = config.Authentication.Audiences
	}
	if config.Authentication.LoginEndpoint != "" {
		e.Authorization.LoginEndpoint = config.Authentication.LoginEndpoint
	}
	if config.Authentication.IdentityProvider != "" {
		e.Authorization.IdentityProvider = config.Authentication.IdentityProvider
	}
	if config.Authentication.Tenant != "" {
		e.Authorization.Tenant = config.Authentication.Tenant
	}

	if config.DnsSuffixes.Attestation != "" && config.ResourceIdentifiers.Attestation != "" {
		e.Attestation = AttestationAPI(config.ResourceIdentifiers.Attestation, config.DnsSuffixes.Attestation)
	}
	if config.ResourceIdentifiers.Batch != "" {
		e.Batch = BatchAPI(config.ResourceIdentifiers.Batch)
	}
	if config.DnsSuffixes.FrontDoor != "" {
		e.CDNFrontDoor = CDNFrontDoorAPI(config.DnsSuffixes.FrontDoor)
	}
	if config.DnsSuffixes.ContainerRegistry != "" {
		e.ContainerRegistry = ContainerRegistryAPI(config.DnsSuffixes.ContainerRegistry)
	}
	if config.DnsSuffixes.DataLakeStore != "" && config.ResourceIdentifiers.DataLake != "" {
		e.DataLake = DataLakeAPI(config.DnsSuffixes.DataLakeStore).WithResourceIdentifier(config.ResourceIdentifiers.DataLake)
	}
	if config.DnsSuffixes.KeyVault != "" {
		// Key Vault resource ID is missing in metadata, so make a best-effort guess from the domain suffix
		e.KeyVault = KeyVaultAPI(config.DnsSuffixes.KeyVault).WithResourceIdentifier(fmt.Sprintf("https://%s", config.DnsSuffixes.KeyVault))
	}
	if config.DnsSuffixes.ManagedHSM != "" {
		// Managed HSM resource ID is missing in metadata, so make a best-effort guess from the domain suffix
		mHsmEndpoint := fmt.Sprintf("https://%s", config.DnsSuffixes.ManagedHSM)
		e.ManagedHSM = ManagedHSMAPI(mHsmEndpoint, config.DnsSuffixes.ManagedHSM).WithResourceIdentifier(mHsmEndpoint)
	}
	if config.DnsSuffixes.MariaDB != "" && config.ResourceIdentifiers.OSSRDBMS != "" {
		e.MariaDB = MariaDBAPI(config.DnsSuffixes.MariaDB).WithResourceIdentifier(config.ResourceIdentifiers.OSSRDBMS)
	}
	if config.ResourceIdentifiers.MicrosoftGraph != "" {
		e.MicrosoftGraph = MicrosoftGraphAPI(config.ResourceIdentifiers.MicrosoftGraph)
	}
	if config.DnsSuffixes.MySql != "" && config.ResourceIdentifiers.OSSRDBMS != "" {
		e.MySql = MySqlAPI(config.DnsSuffixes.MySql).WithResourceIdentifier(config.ResourceIdentifiers.OSSRDBMS)
	}
	if config.ResourceIdentifiers.LogAnalytics != "" {
		e.OperationalInsights = OperationalInsightsAPI().WithResourceIdentifier(config.ResourceIdentifiers.LogAnalytics)
	}
	if config.DnsSuffixes.Postgresql != "" && config.ResourceIdentifiers.OSSRDBMS != "" {
		e.Postgresql = PostgresqlAPI(config.DnsSuffixes.Postgresql).WithResourceIdentifier(config.ResourceIdentifiers.OSSRDBMS)
	}
	if config.ResourceManagerEndpoint != "" {
		e.ResourceManager = ResourceManagerAPI(config.ResourceManagerEndpoint)
	}
	if config.DnsSuffixes.SqlServer != "" {
		// SQL resource ID is missing in metadata, so make a best-effort guess from the domain suffix
		e.Sql = SqlAPI(config.DnsSuffixes.SqlServer).WithResourceIdentifier(fmt.Sprintf("https://%s", config.DnsSuffixes.SqlServer))
	}
	if config.DnsSuffixes.Storage != "" {
		e.Storage = StorageAPI(config.DnsSuffixes.Storage)
	}
	if config.DnsSuffixes.StorageSync != "" {
		e.StorageSync = StorageSyncAPI(config.DnsSuffixes.StorageSync)
	}
	if config.DnsSuffixes.Synapse != "" && config.ResourceIdentifiers.Synapse != "" {
		e.Synapse = SynapseAPI(config.DnsSuffixes.Synapse).WithResourceIdentifier(config.ResourceIdentifiers.Synapse)
	}

	return nil
}
