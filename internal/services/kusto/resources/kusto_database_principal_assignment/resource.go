// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto_database_principal_assignment

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/databaseprincipalassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/resources/kusto_database_principal_assignment/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKustoDatabasePrincipalAssignmentCreate,
		Read:   resourceKustoDatabasePrincipalAssignmentRead,
		Delete: resourceKustoDatabasePrincipalAssignmentDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.KustoDatabasePrincipalAssignmentV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := databaseprincipalassignments.ParseDatabasePrincipalAssignmentID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: kustoDatabasePrincipalAssignmentSchema(),
	}
}
