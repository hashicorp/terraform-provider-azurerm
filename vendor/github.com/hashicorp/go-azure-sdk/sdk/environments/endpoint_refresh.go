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
	config, err := client.GetMetaData(ctx, e.Name)
	if err != nil {
		return fmt.Errorf("retrieving MetaData from endpoint: %+v", err)
	}

	if err := e.updateFromMetaData(config); err != nil {
		return fmt.Errorf("updating Environment from MetaData: %+v", err)
	}

	return nil
}

func (e *Environment) updateFromMetaData(config *metadata.MetaData) error {
	// Auth Endpoints
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
	if config.ResourceManagerEndpoint != "" {
		e.ResourceManager = ResourceManagerAPI(config.ResourceManagerEndpoint)
	}
	if config.ResourceIdentifiers.MicrosoftGraph != "" {
		e.MicrosoftGraph = MicrosoftGraphAPI(config.ResourceIdentifiers.MicrosoftGraph)
	}

	// Dns Suffixes
	if config.DnsSuffixes.FrontDoor != "" {
		e.CDNFrontDoor = CDNFrontDoorAPI(config.DnsSuffixes.FrontDoor)
	}
	if config.DnsSuffixes.KeyVault != "" {
		e.KeyVault = KeyVaultAPI(config.DnsSuffixes.KeyVault)
	}
	if config.DnsSuffixes.ManagedHSM != "" {
		e.ManagedHSM = ManagedHSMAPI(fmt.Sprintf("https://%s", config.DnsSuffixes.ManagedHSM), config.DnsSuffixes.ManagedHSM)
	}
	if config.DnsSuffixes.MariaDB != "" {
		e.MariaDB = MariaDBAPI(config.DnsSuffixes.MariaDB)
	}
	if config.DnsSuffixes.MySql != "" {
		e.MySql = MySqlAPI(config.DnsSuffixes.MySql)
	}
	if config.DnsSuffixes.Postgresql != "" {
		e.Postgresql = PostgresqlAPI(config.DnsSuffixes.Postgresql)
	}
	if config.DnsSuffixes.SqlServer != "" {
		e.Sql = SqlAPI(config.DnsSuffixes.SqlServer)
	}
	if config.DnsSuffixes.Storage != "" {
		e.Storage = StorageAPI(config.DnsSuffixes.Storage)
	}
	if config.DnsSuffixes.StorageSync != "" {
		e.StorageSync = StorageSyncAPI(config.DnsSuffixes.StorageSync)
	}
	if config.DnsSuffixes.Synapse != "" {
		e.Synapse = SynapseAPI(config.DnsSuffixes.Synapse)
	}

	return nil
}
