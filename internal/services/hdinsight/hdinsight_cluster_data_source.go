// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hdinsight

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/hdinsight/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceHDInsightSparkCluster() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceHDInsightClusterRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": SchemaHDInsightDataSourceName(),

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

			"tags": tags.SchemaDataSource(),

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
		},
	}
}

func dataSourceHDInsightClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	clustersClient := meta.(*clients.Client).HDInsight.ClustersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	configurationsClient := meta.(*clients.Client).HDInsight.ConfigurationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewClusterID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := clustersClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	configuration, err := configurationsClient.Get(ctx, id.ResourceGroup, id.Name, "gateway")
	if err != nil {
		return fmt.Errorf("retrieving Configuration for %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.Properties; props != nil {
		d.Set("cluster_version", props.ClusterVersion)
		d.Set("tier", string(props.Tier))
		d.Set("tls_min_version", props.MinSupportedTLSVersion)

		if def := props.ClusterDefinition; def != nil {
			d.Set("component_versions", flattenHDInsightsDataSourceComponentVersions(def.ComponentVersion))
			if kind := def.Kind; kind != nil {
				d.Set("kind", strings.ToLower(*kind))
			}
			if err := d.Set("gateway", FlattenHDInsightsConfigurations(configuration.Value, d)); err != nil {
				return fmt.Errorf("flattening `gateway`: %+v", err)
			}
		}

		edgeNodeSshEndpoint := FindHDInsightConnectivityEndpoint("EDGESSH", props.ConnectivityEndpoints)
		d.Set("edge_ssh_endpoint", edgeNodeSshEndpoint)
		httpEndpoint := FindHDInsightConnectivityEndpoint("HTTPS", props.ConnectivityEndpoints)
		d.Set("https_endpoint", httpEndpoint)
		sshEndpoint := FindHDInsightConnectivityEndpoint("SSH", props.ConnectivityEndpoints)
		d.Set("ssh_endpoint", sshEndpoint)
		kafkaRestProxyEndpoint := FindHDInsightConnectivityEndpoint("KafkaRestProxyPublicEndpoint", props.ConnectivityEndpoints)
		d.Set("kafka_rest_proxy_endpoint", kafkaRestProxyEndpoint)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func flattenHDInsightsDataSourceComponentVersions(input map[string]*string) map[string]string {
	output := make(map[string]string)

	for k, v := range input {
		if v == nil {
			continue
		}

		output[k] = *v
	}

	return output
}
