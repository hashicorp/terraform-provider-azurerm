// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_cassandra_datacenter

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2023-04-15/managedcassandras"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCassandraDatacenterRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cosmos.ManagedCassandraClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := managedcassandras.ParseDataCenterID(d.Id())
	if err != nil {
		return err
	}
	resp, err := client.CassandraDataCentersGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Error reading %q - removing from state", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading %q: %+v", id, err)
	}

	clusterId := managedcassandras.NewCassandraClusterID(id.SubscriptionId, id.ResourceGroupName, id.CassandraClusterName)
	d.Set("name", id.DataCenterName)
	d.Set("cassandra_cluster_id", clusterId.ID())
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("delegated_management_subnet_id", props.DelegatedSubnetId)
			d.Set("location", location.NormalizeNilable(props.DataCenterLocation))
			d.Set("backup_storage_customer_key_uri", props.BackupStorageCustomerKeyUri)
			d.Set("base64_encoded_yaml_fragment", props.Base64EncodedCassandraYamlFragment)
			d.Set("managed_disk_customer_key_uri", props.ManagedDiskCustomerKeyUri)
			d.Set("node_count", props.NodeCount)
			d.Set("disk_count", int(*props.DiskCapacity))
			d.Set("disk_sku", props.DiskSku)
			d.Set("sku_name", props.Sku)
			d.Set("availability_zones_enabled", props.AvailabilityZone)

			if err := d.Set("seed_node_ip_addresses", flattenCassandraDatacenterSeedNodes(props.SeedNodes)); err != nil {
				return fmt.Errorf("setting `seed_node_ip_addresses`: %+v", err)
			}
		}
	}
	return nil
}

func flattenCassandraDatacenterSeedNodes(input *[]managedcassandras.SeedNode) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		if item.IPAddress != nil {
			results = append(results, item.IPAddress)
		}
	}

	return results
}

func cassandraDatacenterStateRefreshFunc(ctx context.Context, client *managedcassandras.ManagedCassandrasClient, id managedcassandras.DataCenterId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		res, err := client.CassandraDataCentersGet(ctx, id)
		if err != nil {
			return nil, "", fmt.Errorf("polling for %s: %+v", id, err)
		}

		if model := res.Model; model != nil {
			if model.Properties != nil && model.Properties.ProvisioningState != nil {
				return res, string(*model.Properties.ProvisioningState), nil
			}
		}
		return nil, "", fmt.Errorf("unable to read provisioning state")
	}
}
