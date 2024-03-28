// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cosmos

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2023-04-15/managedcassandras"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCassandraDatacenter() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceCassandraDatacenterCreate,
		Read:   resourceCassandraDatacenterRead,
		Update: resourceCassandraDatacenterUpdate,
		Delete: resourceCassandraDatacenterDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := managedcassandras.ParseDataCenterID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"cassandra_cluster_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CassandraClusterID,
			},

			"location": commonschema.Location(),

			"delegated_management_subnet_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateSubnetID,
			},

			"backup_storage_customer_key_uri": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: keyVaultValidate.NestedItemId,
			},

			"base64_encoded_yaml_fragment": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"disk_sku": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      "P30",
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"managed_disk_customer_key_uri": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: keyVaultValidate.NestedItemId,
			},
			"node_count": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(3),
				Default:      3,
			},
			"disk_count": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 10),
			},
			"availability_zones_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"seed_node_ip_addresses": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
		},
	}

	if !features.FourPointOhBeta() {
		// NOTE: The API does not expose a constant for the Sku so I had to hardcode it here...
		// Per the service team, the current default Sku is 'Standard_DS14_v2' but moving forward
		// the new default value should be 'Standard_E16s_v5'.
		resource.Schema["sku_name"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		}
	} else {
		resource.Schema["sku_name"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "Standard_E16s_v5",
			ValidateFunc: validation.StringIsNotEmpty,
		}
	}

	return resource
}

func resourceCassandraDatacenterCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.ManagedCassandraClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	clusterId, err := managedcassandras.ParseCassandraClusterID(d.Get("cassandra_cluster_id").(string))
	if err != nil {
		return err
	}
	id := managedcassandras.NewDataCenterID(clusterId.SubscriptionId, clusterId.ResourceGroupName, clusterId.CassandraClusterName, d.Get("name").(string))

	existing, err := client.CassandraDataCentersGet(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_cosmosdb_cassandra_datacenter", id.ID())
	}

	payload := managedcassandras.DataCenterResource{
		Properties: &managedcassandras.DataCenterResourceProperties{
			DelegatedSubnetId:  utils.String(d.Get("delegated_management_subnet_id").(string)),
			NodeCount:          utils.Int64(int64(d.Get("node_count").(int))),
			AvailabilityZone:   utils.Bool(d.Get("availability_zones_enabled").(bool)),
			DiskCapacity:       utils.Int64(int64(d.Get("disk_count").(int))),
			DiskSku:            utils.String(d.Get("disk_sku").(string)),
			DataCenterLocation: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		},
	}

	if v, ok := d.GetOk("backup_storage_customer_key_uri"); ok {
		payload.Properties.BackupStorageCustomerKeyUri = utils.String(v.(string))
	}

	if v, ok := d.GetOk("base64_encoded_yaml_fragment"); ok {
		payload.Properties.Base64EncodedCassandraYamlFragment = utils.String(v.(string))
	}

	if v, ok := d.GetOk("managed_disk_customer_key_uri"); ok {
		payload.Properties.ManagedDiskCustomerKeyUri = utils.String(v.(string))
	}

	if v, ok := d.GetOk("sku_name"); ok {
		payload.Properties.Sku = utils.String(v.(string))
	}

	if err = client.CassandraDataCentersCreateUpdateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceCassandraDatacenterRead(d, meta)
}

func resourceCassandraDatacenterRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

func resourceCassandraDatacenterUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.ManagedCassandraClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := managedcassandras.ParseDataCenterID(d.Id())
	if err != nil {
		return err
	}

	payload := managedcassandras.DataCenterResource{
		Properties: &managedcassandras.DataCenterResourceProperties{
			DelegatedSubnetId:  utils.String(d.Get("delegated_management_subnet_id").(string)),
			NodeCount:          utils.Int64(int64(d.Get("node_count").(int))),
			Sku:                utils.String(d.Get("sku_name").(string)),
			DataCenterLocation: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
			DiskSku:            utils.String(d.Get("disk_sku").(string)),
		},
	}

	if v, ok := d.GetOk("backup_storage_customer_key_uri"); ok {
		payload.Properties.BackupStorageCustomerKeyUri = utils.String(v.(string))
	}

	if v, ok := d.GetOk("base64_encoded_yaml_fragment"); ok {
		payload.Properties.Base64EncodedCassandraYamlFragment = utils.String(v.(string))
	}

	if v, ok := d.GetOk("managed_disk_customer_key_uri"); ok {
		payload.Properties.ManagedDiskCustomerKeyUri = utils.String(v.(string))
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

func resourceCassandraDatacenterDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.ManagedCassandraClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := managedcassandras.ParseDataCenterID(d.Id())
	if err != nil {
		return err
	}

	if err := client.CassandraDataCentersDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %q: %+v", id, err)
	}

	return nil
}

func cassandraDatacenterStateRefreshFunc(ctx context.Context, client *managedcassandras.ManagedCassandrasClient, id managedcassandras.DataCenterId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
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

func flattenCassandraDatacenterSeedNodes(input *[]managedcassandras.SeedNode) []interface{} {
	results := make([]interface{}, 0)
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
