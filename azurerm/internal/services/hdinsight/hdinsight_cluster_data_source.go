package hdinsight

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceHDInsightSparkCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceHDInsightClusterRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": SchemaHDInsightDataSourceName(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"cluster_version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"component_versions": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"kind": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tier": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tls_min_version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"gateway": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"password": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},

			"tags": tags.SchemaDataSource(),

			"edge_ssh_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"https_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"kafka_rest_proxy_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ssh_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceHDInsightClusterRead(d *schema.ResourceData, meta interface{}) error {
	clustersClient := meta.(*clients.Client).HDInsight.ClustersClient
	configurationsClient := meta.(*clients.Client).HDInsight.ConfigurationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	resp, err := clustersClient.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("HDInsight Cluster %q was not found in Resource Group %q", name, resourceGroup)
		}

		return fmt.Errorf("Error retrieving HDInsight Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	configuration, err := configurationsClient.Get(ctx, resourceGroup, name, "gateway")
	if err != nil {
		return fmt.Errorf("Error retrieving Configuration for HDInsight Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.Properties; props != nil {
		d.Set("cluster_version", props.ClusterVersion)
		d.Set("tier", string(props.Tier))
		d.Set("tls_min_version", props.MinSupportedTLSVersion)

		if def := props.ClusterDefinition; def != nil {
			d.Set("component_versions", flattenHDInsightsDataSourceComponentVersions(def.ComponentVersion))
			if kind := def.Kind; kind != nil {
				d.Set("kind", strings.ToLower(*kind))
			}
			if err := d.Set("gateway", FlattenHDInsightsConfigurations(configuration.Value)); err != nil {
				return fmt.Errorf("Error flattening `gateway`: %+v", err)
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
