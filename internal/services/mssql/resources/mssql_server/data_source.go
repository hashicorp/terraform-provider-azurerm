// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_server

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/sqlvulnerabilityassessmentssettings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Microsoft SQL Server.",
		Read:        dataSourceMsSqlServerRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description: "The name of this Microsoft SQL Server.",
				Type:        pluginsdk.TypeString,
				Required:    true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"version": {
				Description: "The version of the Microsoft SQL Server.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"administrator_login": {
				Description: "The administrator login name of the Microsoft SQL Server.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"express_vulnerability_assessment_enabled": {
				Description: "Whether the Express Vulnerability Assessment Configuration is enabled.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"fully_qualified_domain_name": {
				Description: "The fully qualified domain name of the Microsoft SQL Server.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"identity": commonschema.SystemOrUserAssignedIdentityComputed(),

			"restorable_dropped_database_ids": {
				Description: "A list of dropped restorable database IDs on the Microsoft SQL Server.",
				Type:        pluginsdk.TypeList,
				Computed:    true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"soft_delete_retention_days": {
				Description: "The number of days a deleted server is retained in a soft-deleted state and can be restored. A value of `0` means soft delete is disabled.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"transparent_data_encryption_key_vault_key_id": {
				Description: "The Key Vault Key URI to be used as the `Customer Managed Key` (CMK/BYOK) for the `Transparent Data Encryption` (TDE) layer.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceMsSqlServerRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MSSQL.ServersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	restorableDroppedDatabasesClient := meta.(*clients.Client).MSSQL.RestorableDroppedDatabasesClient
	vaClient := meta.(*clients.Client).MSSQL.SqlVulnerabilityAssessmentSettingsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewSqlServerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id, servers.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %s", id, err)
	}

	d.SetId(id.ID())
	if model := resp.Model; model != nil {
		// NOTE: In the new API Location is just a string, not a pointer..
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			d.Set("version", props.Version)
			d.Set("administrator_login", props.AdministratorLogin)
			d.Set("fully_qualified_domain_name", props.FullyQualifiedDomainName)
			d.Set("soft_delete_retention_days", int(pointer.From(props.RetentionDays)))
			d.Set("transparent_data_encryption_key_vault_key_id", props.KeyId)
		}

		identity, err := identity.FlattenLegacySystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if err := d.Set("identity", identity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		restorableListPage, err := restorableDroppedDatabasesClient.ListByServerComplete(ctx, id)
		if err != nil {
			return fmt.Errorf("listing %s Restorable Dropped Databases: %v", id, err)
		}

		if err := d.Set("restorable_dropped_database_ids", flattenSqlServerRestorableDatabases(restorableListPage)); err != nil {
			return fmt.Errorf("setting `restorable_dropped_database_ids`: %+v", err)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	va, err := vaClient.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving SQL Server Vulnerability Assessment Settings for %s: %+v", id, err)
	}

	if model := va.Model; model != nil && model.Properties != nil {
		d.Set("express_vulnerability_assessment_enabled", pointer.From(model.Properties.State) == sqlvulnerabilityassessmentssettings.SqlVulnerabilityAssessmentStateEnabled)
	}

	return nil
}
