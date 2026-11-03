// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto_attached_database_configuration

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/attacheddatabaseconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKustoAttachedDatabaseConfigurationRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Kusto.AttachedDatabaseConfigurationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := attacheddatabaseconfigurations.ParseAttachedDatabaseConfigurationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.AttachedDatabaseConfigurationName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("cluster_name", id.ClusterName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			clusterResourceId, parseErr := commonids.ParseKustoClusterIDInsensitively(props.ClusterResourceId)
			if parseErr != nil {
				return parseErr
			}
			d.Set("cluster_id", clusterResourceId.ID())
			d.Set("database_name", props.DatabaseName)
			d.Set("database_name_override", pointer.From(props.DatabaseNameOverride))
			d.Set("database_name_prefix", pointer.From(props.DatabaseNamePrefix))
			d.Set("default_principal_modification_kind", props.DefaultPrincipalsModificationKind)
			d.Set("attached_database_names", props.AttachedDatabaseNames)
			d.Set("sharing", flattenAttachedDatabaseConfigurationTableLevelSharingProperties(props.TableLevelSharingProperties))
		}
	}

	return nil
}

func flattenAttachedDatabaseConfigurationTableLevelSharingProperties(input *attacheddatabaseconfigurations.TableLevelSharingProperties) []any {
	if input == nil {
		return make([]any, 0)
	}

	return []any{
		map[string]any{
			"external_tables_to_exclude":    pluginsdk.FlattenSlice(input.ExternalTablesToExclude),
			"external_tables_to_include":    pluginsdk.FlattenSlice(input.ExternalTablesToInclude),
			"functions_to_exclude":          pluginsdk.FlattenSlice(input.FunctionsToExclude),
			"functions_to_include":          pluginsdk.FlattenSlice(input.FunctionsToInclude),
			"materialized_views_to_exclude": pluginsdk.FlattenSlice(input.MaterializedViewsToExclude),
			"materialized_views_to_include": pluginsdk.FlattenSlice(input.MaterializedViewsToInclude),
			"tables_to_exclude":             pluginsdk.FlattenSlice(input.TablesToExclude),
			"tables_to_include":             pluginsdk.FlattenSlice(input.TablesToInclude),
		},
	}
}
