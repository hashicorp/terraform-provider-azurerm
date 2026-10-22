// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_database_extended_auditing_policy

import (
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMsSqlDatabaseExtendedAuditingPolicyCreateUpdate,
		Read:   resourceMsSqlDatabaseExtendedAuditingPolicyRead,
		Update: resourceMsSqlDatabaseExtendedAuditingPolicyCreateUpdate,
		Delete: resourceMsSqlDatabaseExtendedAuditingPolicyDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.DatabaseExtendedAuditingPolicyID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: mssqlDatabaseExtendedAuditingPolicySchema(),
	}
}
