// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package hdinsight_interactive_query_cluster

import (
	"fmt"
	"maps"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/hdinsight/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceHDInsightInteractiveQueryClusterCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).HDInsight.Clusters
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	extensionsClient := meta.(*clients.Client).HDInsight.Extensions
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	id := commonids.NewHDInsightClusterID(subscriptionId, resourceGroup, name)
	location := location.Normalize(d.Get("location").(string))
	clusterVersion := d.Get("cluster_version").(string)
	t := d.Get("tags").(map[string]any)
	tier := clusters.Tier(d.Get("tier").(string))
	tls := d.Get("tls_min_version").(string)

	componentVersionsRaw := d.Get("component_version").([]any)
	componentVersions := expandHDInsightInteractiveQueryComponentVersion(componentVersionsRaw)

	gatewayRaw := d.Get("gateway").([]any)
	configurations := helpers.ExpandHDInsightsConfigurations(gatewayRaw)

	metastoresRaw := d.Get("metastores").([]any)
	metastores := helpers.ExpandHDInsightsMetastore(metastoresRaw)
	maps.Copy(configurations, metastores)

	networkPropertiesRaw := d.Get("network").([]any)
	networkProperties := helpers.ExpandHDInsightsNetwork(networkPropertiesRaw)

	privateLinkConfigurationsRaw := d.Get("private_link_configuration").([]any)
	privateLinkConfigurations := helpers.ExpandHDInsightPrivateLinkConfigurations(privateLinkConfigurationsRaw)

	storageAccountsRaw := d.Get("storage_account").([]any)
	storageAccountsGen2Raw := d.Get("storage_account_gen2").([]any)
	storageAccounts, expandedIdentity, err := helpers.ExpandHDInsightsStorageAccounts(storageAccountsRaw, storageAccountsGen2Raw)
	if err != nil {
		return fmt.Errorf("expanding `storage_account`: %s", err)
	}

	interactiveQueryRoles := helpers.HdInsightRoleDefinition{
		HeadNodeDef:      hdInsightInteractiveQueryClusterHeadNodeDefinition,
		WorkerNodeDef:    hdInsightInteractiveQueryClusterWorkerNodeDefinition,
		ZookeeperNodeDef: hdInsightInteractiveQueryClusterZookeeperNodeDefinition,
	}
	rolesRaw := d.Get("roles").([]any)
	roles, err := helpers.ExpandHDInsightRoles(rolesRaw, interactiveQueryRoles)
	if err != nil {
		return fmt.Errorf("expanding `roles`: %+v", err)
	}

	computeIsolationProperties := helpers.ExpandHDInsightComputeIsolationProperties(d.Get("compute_isolation").([]any))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing HDInsight InteractiveQuery Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_hdinsight_interactive_query_cluster", id.ID())
		}
	}

	encryptionInTransit := d.Get("encryption_in_transit_enabled").(bool)

	var configurationsRaw any = configurations
	params := clusters.ClusterCreateParametersExtended{
		Location: pointer.To(location),
		Properties: &clusters.ClusterCreateProperties{
			Tier:                      pointer.To(tier),
			OsType:                    pointer.To(clusters.OSTypeLinux),
			ClusterVersion:            pointer.To(clusterVersion),
			MinSupportedTlsVersion:    pointer.To(tls),
			NetworkProperties:         networkProperties,
			PrivateLinkConfigurations: privateLinkConfigurations,
			EncryptionInTransitProperties: &clusters.EncryptionInTransitProperties{
				IsEncryptionInTransitEnabled: &encryptionInTransit,
			},
			ClusterDefinition: &clusters.ClusterDefinition{
				Kind:             pointer.To(clusters.ClusterKindInteractiveHive),
				ComponentVersion: componentVersions,
				Configurations:   &configurationsRaw,
			},
			StorageProfile: &clusters.StorageProfile{
				Storageaccounts: storageAccounts,
			},
			ComputeProfile: &clusters.ComputeProfile{
				Roles: roles,
			},
			ComputeIsolationProperties: computeIsolationProperties,
		},
		Tags:     tags.Expand(t),
		Identity: expandedIdentity,
	}

	if diskEncryptionPropertiesRaw, ok := d.GetOk("disk_encryption"); ok {
		params.Properties.DiskEncryptionProperties, err = helpers.ExpandHDInsightsDiskEncryptionProperties(diskEncryptionPropertiesRaw.([]any))
		if err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("security_profile"); ok {
		params.Properties.SecurityProfile = helpers.ExpandHDInsightSecurityProfile(v.([]any))

		// @tombuildsstuff: this behaviour is likely wrong and wants reevaluating - users should need to explicitly define this in the config?
		params.Identity = &identity.SystemAndUserAssignedMap{
			Type:        identity.TypeUserAssigned,
			IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
		}
		if params.Properties.SecurityProfile != nil && params.Properties.SecurityProfile.MsiResourceId != nil {
			params.Identity.IdentityIds[*params.Properties.SecurityProfile.MsiResourceId] = identity.UserAssignedIdentityDetails{
				// intentionally empty
			}
		}
	}

	if err := client.CreateCallbackThenPoll(ctx, id, params, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating Interactive Query %s: %+v", id, err)
	}

	d.SetId(id.ID())

	// We can only enable monitoring after creation
	if v, ok := d.GetOk("monitor"); ok {
		monitorRaw := v.([]any)
		if err := helpers.EnableHDInsightMonitoring(ctx, extensionsClient, id, monitorRaw); err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("extension"); ok {
		extensionRaw := v.([]any)
		if err := helpers.EnableHDInsightAzureMonitor(ctx, extensionsClient, id, extensionRaw); err != nil {
			return err
		}
	}

	return resourceHDInsightInteractiveQueryClusterRead(d, meta)
}

func expandHDInsightInteractiveQueryComponentVersion(input []any) *map[string]string {
	vs := input[0].(map[string]any)
	return &map[string]string{
		"InteractiveHive": vs["interactive_hive"].(string),
	}
}
