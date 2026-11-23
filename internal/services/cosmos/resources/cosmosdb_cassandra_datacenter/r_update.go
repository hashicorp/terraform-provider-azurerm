// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_cassandra_datacenter

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2023-04-15/managedcassandras"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCassandraDatacenterUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cosmos.ManagedCassandraClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := managedcassandras.ParseDataCenterID(d.Id())
	if err != nil {
		return err
	}

	payload := managedcassandras.DataCenterResource{
		Properties: &managedcassandras.DataCenterResourceProperties{
			DelegatedSubnetId:  pointer.To(d.Get("delegated_management_subnet_id").(string)),
			NodeCount:          pointer.To(int64(d.Get("node_count").(int))),
			Sku:                pointer.To(d.Get("sku_name").(string)),
			DataCenterLocation: pointer.To(location.Normalize(d.Get("location").(string))),
			DiskSku:            pointer.To(d.Get("disk_sku").(string)),
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

	if err := client.CassandraDataCentersCreateUpdateThenPoll(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %q: %+v", id, err)
	}

	// Issue: https://github.com/Azure/azure-rest-api-specs/issues/19078
	// There is a long running issue on updating this resource.
	// The API cannot update the property after WaitForCompletionRef is returned.
	// It has to wait a while after that. Then the property can be updated successfully.
	stateConf := &pluginsdk.StateChangeConf{
		Delay:      1 * time.Minute,
		Pending:    []string{string(managedcassandras.ManagedCassandraProvisioningStateUpdating)},
		Target:     []string{string(managedcassandras.ManagedCassandraProvisioningStateSucceeded)},
		Refresh:    cassandraDatacenterStateRefreshFunc(ctx, client, *id),
		MinTimeout: 15 * time.Second,
		Timeout:    d.Timeout(pluginsdk.TimeoutUpdate),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for update of %s: %+v", id, err)
	}

	return resourceCassandraDatacenterRead(d, meta)
}
