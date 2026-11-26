// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package hdinsight_hadoop_cluster

import (
	"fmt"
	"log"
	"maps"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/applications"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/clusters"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/hdinsight/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/hdinsight/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceHDInsightHadoopClusterCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).HDInsight.Clusters
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	extensionsClient := meta.(*clients.Client).HDInsight.Extensions
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewHDInsightClusterID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	location := location.Normalize(d.Get("location").(string))
	clusterVersion := d.Get("cluster_version").(string)
	t := d.Get("tags").(map[string]any)
	tier := clusters.Tier(d.Get("tier").(string))
	tls := d.Get("tls_min_version").(string)

	componentVersionsRaw := d.Get("component_version").([]any)
	componentVersions := expandHDInsightHadoopComponentVersion(componentVersionsRaw)

	gatewayRaw := d.Get("gateway").([]any)
	configurations := helpers.ExpandHDInsightsConfigurations(gatewayRaw)

	metastoresRaw := d.Get("metastores").([]any)
	metastores := helpers.ExpandHDInsightsMetastore(metastoresRaw)
	maps.Copy(configurations, metastores)

	networkPropertiesRaw := d.Get("network").([]any)
	networkProperties := helpers.ExpandHDInsightsNetwork(networkPropertiesRaw)

	privateLinkConfigurationsRaw := d.Get("private_link_configuration").([]any)
	privateLinkConfigurations := helpers.ExpandHDInsightPrivateLinkConfigurations(privateLinkConfigurationsRaw)

	computeIsolationProperties := helpers.ExpandHDInsightComputeIsolationProperties(d.Get("compute_isolation").([]any))

	storageAccountsRaw := d.Get("storage_account").([]any)
	storageAccountsGen2Raw := d.Get("storage_account_gen2").([]any)
	storageAccounts, expandedIdentity, err := helpers.ExpandHDInsightsStorageAccounts(storageAccountsRaw, storageAccountsGen2Raw)
	if err != nil {
		return fmt.Errorf("expanding `storage_account`: %s", err)
	}

	rolesRaw := d.Get("roles").([]any)
	hadoopRoles := helpers.HdInsightRoleDefinition{
		HeadNodeDef:      hdInsightHadoopClusterHeadNodeDefinition,
		WorkerNodeDef:    hdInsightHadoopClusterWorkerNodeDefinition,
		ZookeeperNodeDef: hdInsightHadoopClusterZookeeperNodeDefinition,
	}
	roles, err := helpers.ExpandHDInsightRoles(rolesRaw, hadoopRoles)
	if err != nil {
		return fmt.Errorf("expanding `roles`: %+v", err)
	}

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Hadoop %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_hdinsight_hadoop_cluster", id.ID())
		}
	}

	var configurationsRaw any = configurations
	payload := clusters.ClusterCreateParametersExtended{
		Location: pointer.To(location),
		Properties: &clusters.ClusterCreateProperties{
			Tier:                      pointer.To(tier),
			OsType:                    pointer.To(clusters.OSTypeLinux),
			ClusterVersion:            pointer.To(clusterVersion),
			MinSupportedTlsVersion:    pointer.To(tls),
			NetworkProperties:         networkProperties,
			PrivateLinkConfigurations: privateLinkConfigurations,
			ClusterDefinition: &clusters.ClusterDefinition{
				Kind:             pointer.To(clusters.ClusterKindHadoop),
				ComponentVersion: pointer.To(componentVersions),
				Configurations:   pointer.To(configurationsRaw),
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
		diskEncryptionProperties, err := helpers.ExpandHDInsightsDiskEncryptionProperties(diskEncryptionPropertiesRaw.([]any))
		if err != nil {
			return err
		}
		payload.Properties.DiskEncryptionProperties = diskEncryptionProperties
	}

	if v, ok := d.GetOk("security_profile"); ok {
		payload.Properties.SecurityProfile = helpers.ExpandHDInsightSecurityProfile(v.([]any))

		// @tombuildsstuff: this behaviour is likely wrong and wants reevaluating - users should need to explicitly define this in the config?
		payload.Identity = &identity.SystemAndUserAssignedMap{
			Type:        identity.TypeUserAssigned,
			IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
		}
		if payload.Properties.SecurityProfile != nil && payload.Properties.SecurityProfile.MsiResourceId != nil {
			payload.Identity.IdentityIds[*payload.Properties.SecurityProfile.MsiResourceId] = identity.UserAssignedIdentityDetails{
				// intentionally empty
			}
		}
	}

	if err := client.CreateCallbackThenPoll(ctx, id, payload, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating Hadoop %s: %+v", id, err)
	}
	d.SetId(id.ID())

	// We can only add an edge node after creation
	if v, ok := d.GetOk("roles.0.edge_node"); ok {
		edgeNodeRaw := v.([]any)
		applicationsClient := meta.(*clients.Client).HDInsight.Applications
		edgeNodeConfig := edgeNodeRaw[0].(map[string]any)
		applicationId := applications.NewApplicationID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName, id.ClusterName) // 2 id.ClusterName's are intentional

		if err := helpers.CreateHDInsightEdgeNodes(ctx, applicationsClient, applicationId, edgeNodeConfig); err != nil {
			return err
		}

		// we can't rely on the use of the Future here due to the node being successfully completed but now the cluster is applying those changes.
		log.Printf("[DEBUG] Waiting for the Hadoop %s to finish applying edge nodes..", id)
		pollerType := custompollers.NewEdgeNodePoller(client, id)
		poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
		if err := poller.PollUntilDone(ctx); err != nil {
			return fmt.Errorf("waiting for the Hadoop %s to finish applying the edge nodes: %+v", id, err)
		}
	}

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

	return resourceHDInsightHadoopClusterRead(d, meta)
}

func expandHDInsightHadoopComponentVersion(input []any) map[string]string {
	vs := input[0].(map[string]any)
	return map[string]string{
		"Hadoop": vs["hadoop"].(string),
	}
}
