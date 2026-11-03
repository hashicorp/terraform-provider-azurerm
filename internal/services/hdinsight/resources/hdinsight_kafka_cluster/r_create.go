// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package hdinsight_kafka_cluster

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

func resourceHDInsightKafkaClusterCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).HDInsight.Clusters
	extensionsClient := meta.(*clients.Client).HDInsight.Extensions

	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewHDInsightClusterID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	location := location.Normalize(d.Get("location").(string))
	clusterVersion := d.Get("cluster_version").(string)
	t := d.Get("tags").(map[string]any)
	tier := clusters.Tier(d.Get("tier").(string))
	tls := d.Get("tls_min_version").(string)

	componentVersionsRaw := d.Get("component_version").([]any)
	componentVersions := expandHDInsightKafkaComponentVersion(componentVersionsRaw)

	gatewayRaw := d.Get("gateway").([]any)
	configurations := helpers.ExpandHDInsightsConfigurations(gatewayRaw)

	metastoresRaw := d.Get("metastores").([]any)
	metastores := helpers.ExpandHDInsightsMetastore(metastoresRaw)
	maps.Copy(configurations, metastores)

	storageAccountsRaw := d.Get("storage_account").([]any)
	storageAccountsGen2Raw := d.Get("storage_account_gen2").([]any)
	storageAccounts, expandedIdentity, err := helpers.ExpandHDInsightsStorageAccounts(storageAccountsRaw, storageAccountsGen2Raw)
	if err != nil {
		return fmt.Errorf("failure expanding `storage_account`: %s", err)
	}

	networkPropertiesRaw := d.Get("network").([]any)
	networkProperties := helpers.ExpandHDInsightsNetwork(networkPropertiesRaw)

	privateLinkConfigurationsRaw := d.Get("private_link_configuration").([]any)
	privateLinkConfigurations := helpers.ExpandHDInsightPrivateLinkConfigurations(privateLinkConfigurationsRaw)

	kafkaRoles := helpers.HdInsightRoleDefinition{
		HeadNodeDef:            hdInsightKafkaClusterHeadNodeDefinition,
		WorkerNodeDef:          hdInsightKafkaClusterWorkerNodeDefinition,
		ZookeeperNodeDef:       hdInsightKafkaClusterZookeeperNodeDefinition,
		KafkaManagementNodeDef: &hdInsightKafkaClusterKafkaManagementNodeDefinition,
	}
	rolesRaw := d.Get("roles").([]any)
	roles, err := helpers.ExpandHDInsightRoles(rolesRaw, kafkaRoles)
	if err != nil {
		return fmt.Errorf("failure expanding `roles`: %+v", err)
	}

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for the presence of an existing Kafka %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_hdinsight_kafka_cluster", id.ID())
		}
	}

	kafkaRestProperty := expandKafkaRestProxyProperty(d.Get("rest_proxy").([]any))

	computeIsolationProperties := helpers.ExpandHDInsightComputeIsolationProperties(d.Get("compute_isolation").([]any))

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
				Kind:             pointer.To(clusters.ClusterKindKafka),
				ComponentVersion: pointer.To(componentVersions),
				Configurations:   pointer.To(configurationsRaw),
			},
			StorageProfile: &clusters.StorageProfile{
				Storageaccounts: storageAccounts,
			},
			ComputeProfile: &clusters.ComputeProfile{
				Roles: roles,
			},
			KafkaRestProperties:        kafkaRestProperty,
			ComputeIsolationProperties: computeIsolationProperties,
		},
		Tags:     tags.Expand(t),
		Identity: expandedIdentity,
	}

	if encryptionInTransit, ok := d.GetOk("encryption_in_transit_enabled"); ok {
		payload.Properties.EncryptionInTransitProperties = &clusters.EncryptionInTransitProperties{
			IsEncryptionInTransitEnabled: pointer.To(encryptionInTransit.(bool)),
		}
	}

	if diskEncryptionPropertiesRaw, ok := d.GetOk("disk_encryption"); ok {
		payload.Properties.DiskEncryptionProperties, err = helpers.ExpandHDInsightsDiskEncryptionProperties(diskEncryptionPropertiesRaw.([]any))
		if err != nil {
			return err
		}
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
		return fmt.Errorf("creating Kafka %s: %+v", id, err)
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

	return resourceHDInsightKafkaClusterRead(d, meta)
}

func expandHDInsightKafkaComponentVersion(input []any) map[string]string {
	kafkaVersion := ""
	if len(input) > 0 && input[0] != nil {
		vs := input[0].(map[string]any)
		kafkaVersion = vs["kafka"].(string)
	}
	return map[string]string{
		"kafka": kafkaVersion,
	}
}

func expandKafkaRestProxyProperty(input []any) *clusters.KafkaRestProperties {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]any)

	return &clusters.KafkaRestProperties{
		ClientGroupInfo: &clusters.ClientGroupInfo{
			GroupId:   pointer.To(raw["security_group_id"].(string)),
			GroupName: pointer.To(raw["security_group_name"].(string)),
		},
	}
}
