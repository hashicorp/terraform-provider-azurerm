// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_cassandra_datacenter

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2023-04-15/managedcassandras"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCassandraDatacenterCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cosmos.ManagedCassandraClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	clusterId, err := managedcassandras.ParseCassandraClusterID(d.Get("cassandra_cluster_id").(string))
	if err != nil {
		return err
	}
	id := managedcassandras.NewDataCenterID(clusterId.SubscriptionId, clusterId.ResourceGroupName, clusterId.CassandraClusterName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.CassandraDataCentersGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_cosmosdb_cassandra_datacenter", id.ID())
		}
	}

	payload := managedcassandras.DataCenterResource{
		Properties: &managedcassandras.DataCenterResourceProperties{
			DelegatedSubnetId:  pointer.To(d.Get("delegated_management_subnet_id").(string)),
			NodeCount:          pointer.To(int64(d.Get("node_count").(int))),
			AvailabilityZone:   pointer.To(d.Get("availability_zones_enabled").(bool)),
			DiskCapacity:       pointer.To(int64(d.Get("disk_count").(int))),
			DiskSku:            pointer.To(d.Get("disk_sku").(string)),
			DataCenterLocation: pointer.To(location.Normalize(d.Get("location").(string))),
		},
	}

	if v, ok := d.GetOk("backup_storage_customer_key_uri"); ok {
		payload.Properties.BackupStorageCustomerKeyUri = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("base64_encoded_yaml_fragment"); ok {
		payload.Properties.Base64EncodedCassandraYamlFragment = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("managed_disk_customer_key_uri"); ok {
		payload.Properties.ManagedDiskCustomerKeyUri = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("sku_name"); ok {
		payload.Properties.Sku = pointer.To(v.(string))
	}

	if err = client.CassandraDataCentersCreateUpdateCallbackThenPoll(ctx, id, payload, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceCassandraDatacenterRead(d, meta)
}
