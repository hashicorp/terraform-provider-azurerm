package cosmos

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-10-15/documentdb"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCassandraDatacenter() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCassandraDatacenterCreate,
		Read:   resourceCassandraDatacenterRead,
		Update: resourceCassandraDatacenterUpdate,
		Delete: resourceCassandraDatacenterDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.CassandraDatacenterID(id)
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
				ValidateFunc: networkValidate.SubnetID,
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
			"sku_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
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
		},
	}
}

func resourceCassandraDatacenterCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CassandraDatacentersClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	clusterId, _ := parse.CassandraClusterID(d.Get("cassandra_cluster_id").(string))
	id := parse.NewCassandraDatacenterID(clusterId.SubscriptionId, clusterId.ResourceGroup, clusterId.Name, d.Get("name").(string))

	existing, err := client.Get(ctx, id.ResourceGroup, id.CassandraClusterName, id.DataCenterName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_cosmosdb_cassandra_datacenter", id.ID())
	}

	body := documentdb.DataCenterResource{
		Properties: &documentdb.DataCenterResourceProperties{
			DelegatedSubnetID:  utils.String(d.Get("delegated_management_subnet_id").(string)),
			NodeCount:          utils.Int32(int32(d.Get("node_count").(int))),
			Sku:                utils.String(d.Get("sku_name").(string)),
			AvailabilityZone:   utils.Bool(d.Get("availability_zones_enabled").(bool)),
			DiskCapacity:       utils.Int32(int32(d.Get("disk_count").(int))),
			DiskSku:            utils.String(d.Get("disk_sku").(string)),
			DataCenterLocation: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		},
	}

	if v, ok := d.GetOk("backup_storage_customer_key_uri"); ok {
		body.Properties.BackupStorageCustomerKeyURI = utils.String(v.(string))
	}

	if v, ok := d.GetOk("base64_encoded_yaml_fragment"); ok {
		body.Properties.Base64EncodedCassandraYamlFragment = utils.String(v.(string))
	}

	if v, ok := d.GetOk("managed_disk_customer_key_uri"); ok {
		body.Properties.ManagedDiskCustomerKeyURI = utils.String(v.(string))
	}

	future, err := client.CreateUpdate(ctx, id.ResourceGroup, id.CassandraClusterName, id.DataCenterName, body)
	if err != nil {
		return fmt.Errorf("creating %q: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creation for %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceCassandraDatacenterRead(d, meta)
}

func resourceCassandraDatacenterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CassandraDatacentersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CassandraDatacenterID(d.Id())
	if err != nil {
		return err
	}
	resp, err := client.Get(ctx, id.ResourceGroup, id.CassandraClusterName, id.DataCenterName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading %q - removing from state", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading %q: %+v", id, err)
	}

	clusterId := parse.NewCassandraClusterID(id.SubscriptionId, id.ResourceGroup, id.CassandraClusterName)
	d.Set("name", id.DataCenterName)
	d.Set("cassandra_cluster_id", clusterId.ID())
	if props := resp.Properties; props != nil {
		if res := props; res != nil {
			d.Set("delegated_management_subnet_id", props.DelegatedSubnetID)
			d.Set("location", location.NormalizeNilable(props.DataCenterLocation))
			d.Set("backup_storage_customer_key_uri", props.BackupStorageCustomerKeyURI)
			d.Set("base64_encoded_yaml_fragment", props.Base64EncodedCassandraYamlFragment)
			d.Set("managed_disk_customer_key_uri", props.ManagedDiskCustomerKeyURI)
			d.Set("node_count", props.NodeCount)
			d.Set("disk_count", int(*props.DiskCapacity))
			d.Set("disk_sku", props.DiskSku)
			d.Set("sku_name", props.Sku)
			d.Set("availability_zones_enabled", props.AvailabilityZone)
		}
	}
	return nil
}

func resourceCassandraDatacenterUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CassandraDatacentersClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CassandraDatacenterID(d.Id())
	if err != nil {
		return err
	}

	body := documentdb.DataCenterResource{
		Properties: &documentdb.DataCenterResourceProperties{
			DelegatedSubnetID:  utils.String(d.Get("delegated_management_subnet_id").(string)),
			NodeCount:          utils.Int32(int32(d.Get("node_count").(int))),
			DataCenterLocation: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
			DiskSku:            utils.String(d.Get("disk_sku").(string)),
		},
	}

	if v, ok := d.GetOk("backup_storage_customer_key_uri"); ok {
		body.Properties.BackupStorageCustomerKeyURI = utils.String(v.(string))
	}

	if v, ok := d.GetOk("base64_encoded_yaml_fragment"); ok {
		body.Properties.Base64EncodedCassandraYamlFragment = utils.String(v.(string))
	}

	if v, ok := d.GetOk("managed_disk_customer_key_uri"); ok {
		body.Properties.ManagedDiskCustomerKeyURI = utils.String(v.(string))
	}

	future, err := client.CreateUpdate(ctx, id.ResourceGroup, id.CassandraClusterName, id.DataCenterName, body)
	if err != nil {
		return fmt.Errorf("updating %q: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on update for %q: %+v", id, err)
	}

	// Issue: https://github.com/Azure/azure-rest-api-specs/issues/19078
	// There is a long running issue on updating this resource.
	// The API cannot update the property after WaitForCompletionRef is returned.
	// It has to wait a while after that. Then the property can be updated successfully.
	stateConf := &pluginsdk.StateChangeConf{
		Delay:      1 * time.Minute,
		Pending:    []string{string(documentdb.ManagedCassandraProvisioningStateUpdating)},
		Target:     []string{string(documentdb.ManagedCassandraProvisioningStateSucceeded)},
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
	client := meta.(*clients.Client).Cosmos.CassandraDatacentersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CassandraDatacenterID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.CassandraClusterName, id.DataCenterName)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("deleting %q: %+v", id, err)
		}
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("waiting on deleting for %q: %+v", id, err)
	}

	return nil
}

func cassandraDatacenterStateRefreshFunc(ctx context.Context, client *documentdb.CassandraDataCentersClient, id parse.CassandraDatacenterId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id.ResourceGroup, id.CassandraClusterName, id.DataCenterName)
		if err != nil {
			return nil, "", fmt.Errorf("polling for %s: %+v", id, err)
		}

		if res.Properties != nil && res.Properties.ProvisioningState != "" {
			return res, string(res.Properties.ProvisioningState), nil
		}
		return nil, "", fmt.Errorf("unable to read provisioning state")
	}
}
