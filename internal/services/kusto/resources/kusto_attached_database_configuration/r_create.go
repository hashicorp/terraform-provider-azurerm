// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto_attached_database_configuration

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/attacheddatabaseconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKustoAttachedDatabaseConfigurationCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Kusto.AttachedDatabaseConfigurationsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := attacheddatabaseconfigurations.NewAttachedDatabaseConfigurationID(subscriptionId, d.Get("resource_group_name").(string), d.Get("cluster_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			resp, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(resp.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_kusto_attached_database_configuration", id.ID())
			}
		}
	}

	configurationProperties := expandKustoAttachedDatabaseConfigurationProperties(d)
	configurationRequest := attacheddatabaseconfigurations.AttachedDatabaseConfiguration{
		Location:   pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: configurationProperties,
	}

	if d.IsNewResource() {
		if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, configurationRequest, sdk.SetIDCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
		d.SetId(id.ID())
	} else {
		if err := client.CreateOrUpdateThenPoll(ctx, id, configurationRequest); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	return resourceKustoAttachedDatabaseConfigurationRead(d, meta)
}

func expandKustoAttachedDatabaseConfigurationProperties(d *pluginsdk.ResourceData) *attacheddatabaseconfigurations.AttachedDatabaseConfigurationProperties {
	AttachedDatabaseConfigurationProperties := &attacheddatabaseconfigurations.AttachedDatabaseConfigurationProperties{}

	if clusterResourceID, ok := d.GetOk("cluster_id"); ok {
		AttachedDatabaseConfigurationProperties.ClusterResourceId = clusterResourceID.(string)
	}
	if databaseName, ok := d.GetOk("database_name"); ok {
		AttachedDatabaseConfigurationProperties.DatabaseName = databaseName.(string)
	}

	if databaseNameOverride, ok := d.GetOk("database_name_override"); ok {
		AttachedDatabaseConfigurationProperties.DatabaseNameOverride = pointer.To(databaseNameOverride.(string))
	}

	if databaseNamePrefix, ok := d.GetOk("database_name_prefix"); ok {
		AttachedDatabaseConfigurationProperties.DatabaseNamePrefix = pointer.To(databaseNamePrefix.(string))
	}

	if defaultPrincipalModificationKind, ok := d.GetOk("default_principal_modification_kind"); ok {
		AttachedDatabaseConfigurationProperties.DefaultPrincipalsModificationKind = attacheddatabaseconfigurations.DefaultPrincipalsModificationKind(defaultPrincipalModificationKind.(string))
	}

	AttachedDatabaseConfigurationProperties.TableLevelSharingProperties = expandAttachedDatabaseConfigurationTableLevelSharingProperties(d.Get("sharing").([]any))

	return AttachedDatabaseConfigurationProperties
}

func expandAttachedDatabaseConfigurationTableLevelSharingProperties(input []any) *attacheddatabaseconfigurations.TableLevelSharingProperties {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]any)
	return &attacheddatabaseconfigurations.TableLevelSharingProperties{
		TablesToInclude:            pluginsdk.ExpandStringSlice(v["tables_to_include"].(*pluginsdk.Set).List()),
		TablesToExclude:            pluginsdk.ExpandStringSlice(v["tables_to_exclude"].(*pluginsdk.Set).List()),
		ExternalTablesToInclude:    pluginsdk.ExpandStringSlice(v["external_tables_to_include"].(*pluginsdk.Set).List()),
		ExternalTablesToExclude:    pluginsdk.ExpandStringSlice(v["external_tables_to_exclude"].(*pluginsdk.Set).List()),
		FunctionsToInclude:         pluginsdk.ExpandStringSlice(v["functions_to_include"].(*pluginsdk.Set).List()),
		FunctionsToExclude:         pluginsdk.ExpandStringSlice(v["functions_to_exclude"].(*pluginsdk.Set).List()),
		MaterializedViewsToInclude: pluginsdk.ExpandStringSlice(v["materialized_views_to_include"].(*pluginsdk.Set).List()),
		MaterializedViewsToExclude: pluginsdk.ExpandStringSlice(v["materialized_views_to_exclude"].(*pluginsdk.Set).List()),
	}
}
