// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hdinsight

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/configurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/hdinsight/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceHDInsightSparkCluster() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceHDInsightClusterRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.HDInsightName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"cluster_version": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"component_versions": {
				Type:     pluginsdk.TypeMap,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"kind": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tier": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tls_min_version": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"gateway": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
						"username": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"password": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},

			"tags": commonschema.TagsDataSource(),

			"edge_ssh_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"https_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"kafka_rest_proxy_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"ssh_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"cluster_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceHDInsightClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	clustersClient := meta.(*clients.Client).HDInsight.Clusters
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	configurationsClient := meta.(*clients.Client).HDInsight.Configurations
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewHDInsightClusterID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := clustersClient.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	configurationId := configurations.NewConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName, "gateway")
	configurationResp, err := configurationsClient.Get(ctx, configurationId)
	if err != nil {
		return fmt.Errorf("retrieving Configuration for %s: %+v", id, err)
	}

	configuration := make(map[string]string)
	if model := configurationResp.Model; model != nil {
		configuration = *model
	}

	d.SetId(id.ID())

	d.Set("name", id.ClusterName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			d.Set("cluster_id", props.ClusterId)
			d.Set("cluster_version", props.ClusterVersion)
			d.Set("tier", string(pointer.From(props.Tier)))
			d.Set("tls_min_version", props.MinSupportedTlsVersion)

			d.Set("component_versions", flattenHDInsightsDataSourceComponentVersions(props.ClusterDefinition.ComponentVersion))
			d.Set("kind", string(pointer.From(props.ClusterDefinition.Kind)))
			if err := d.Set("gateway", FlattenHDInsightsConfigurations(configuration, d)); err != nil {
				return fmt.Errorf("flattening `gateway`: %+v", err)
			}

			edgeNodeSshEndpoint := findHDInsightConnectivityEndpoint("EDGESSH", props.ConnectivityEndpoints)
			d.Set("edge_ssh_endpoint", edgeNodeSshEndpoint)
			httpEndpoint := findHDInsightConnectivityEndpoint("HTTPS", props.ConnectivityEndpoints)
			d.Set("https_endpoint", httpEndpoint)
			sshEndpoint := findHDInsightConnectivityEndpoint("SSH", props.ConnectivityEndpoints)
			d.Set("ssh_endpoint", sshEndpoint)
			kafkaRestProxyEndpoint := findHDInsightConnectivityEndpoint("KafkaRestProxyPublicEndpoint", props.ConnectivityEndpoints)
			d.Set("kafka_rest_proxy_endpoint", kafkaRestProxyEndpoint)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func flattenHDInsightsDataSourceComponentVersions(input *map[string]string) map[string]string {
	output := make(map[string]string)

	if input != nil {
		for k, v := range *input {
			output[k] = v
		}
	}

	return output
}
