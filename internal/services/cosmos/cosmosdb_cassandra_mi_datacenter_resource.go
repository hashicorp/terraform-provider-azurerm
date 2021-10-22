package cosmos

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-06-15/documentdb"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCassandraMIDatacenter() *pluginsdk.Resource {
	log.Println("manually merged API changes - resourceCassandraMIDatacenter**********")
	return &pluginsdk.Resource{
		Create: resourceCassandraMIDatacenterCreate,
		Read:   resourceCassandraMIDatacenterRead,
		Update: resourceCassandraMIDatacenterUpdate,
		Delete: resourceCassandraMIDatacenterDelete,

		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.CassandraKeyspaceV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"cluster_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CosmosEntityName,
			},
			"datacenter_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CosmosEntityName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			// "location": {
			// 	Type:         pluginsdk.TypeString,
			// 	Required:     true,
			// 	ForceNew:     true,
			// 	ValidateFunc: validate.CosmosAccountName,
			// },
			"location": azure.SchemaLocation(),

			"delegated_management_subnet_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				// ValidateFunc: validate.CosmosEntityName,
				ValidateFunc: networkValidate.SubnetID,
			},

			"node_count": {
				// Type:         pluginsdk.TypeString,
				Type:     pluginsdk.TypeInt,
				Required: true,
				ForceNew: false,
				// ValidateFunc: validate.CosmosEntityName,
			},
		},
	}
}

func resourceCassandraMIDatacenterCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	log.Println("in updated *** resourceCassandraMIDatacenterCreate **********")
	client := meta.(*clients.Client).Cosmos.CassandraDatacentersClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	//name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	clusterName := d.Get("cluster_name").(string)
	datacenterName := d.Get("datacenter_name").(string)
	// nodeCount := d.Get("node_count").(string)
	// nodeCountInt, err2 := strconv.ParseInt(nodeCount, 10, 32)
	// nodeCountInt32 := int32(nodeCountInt)
	location := d.Get("location").(string)
	delegatedSubnetId := d.Get("delegated_management_subnet_id").(string)

	existing, err := client.Get(ctx, resourceGroup, clusterName, datacenterName)

	// if err2 != nil {
	// }
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of creating Cassandra MI  %q (Datacenter: %q): %+v", clusterName, location, err)
		}
	} else {
		if existing.ID == nil && *existing.ID == "" {
			return fmt.Errorf("generating import ID for Cassandra MI  %q (Datacenter: %q)", clusterName, location)
		}

		return tf.ImportAsExistsError("azurerm_cosmosdb_cassandra_mi_datacenter", *existing.ID)
	}

	body := documentdb.DataCenterResource{
		//Location: &location,
		Properties: &documentdb.DataCenterResourceProperties{
			DelegatedSubnetID:  &delegatedSubnetId,
			NodeCount:          utils.Int32(int32(d.Get("node_count").(int))),
			DataCenterLocation: &location,
		},
	}

	future, err := client.CreateUpdate(ctx, resourceGroup, clusterName, datacenterName, body)
	if err != nil {
		return fmt.Errorf("issuing create request for Cassandra MI  %q (Datacenter: %q): %+v", clusterName, datacenterName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on create/update future for Cassandra MI  %q (Datacenter: %q): %+v", clusterName, datacenterName, err)
	}

	resp, err := client.Get(ctx, resourceGroup, clusterName, datacenterName)

	if err != nil {
		return fmt.Errorf("making get request for Cassandra MI  %q (Datacenter: %q): %+v", clusterName, datacenterName, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("getting ID from Cassandra MI  %q (Datacenter: %q)", clusterName, datacenterName)
	}

	d.SetId(*resp.ID)

	return resourceCassandraMIDatacenterRead(d, meta)
}

func resourceCassandraMIDatacenterUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	log.Println("in  *** resourceCassandraMIDatacenterUpdate **********")
	client := meta.(*clients.Client).Cosmos.CassandraDatacentersClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	location := d.Get("location").(string)
	delegatedSubnetId := d.Get("delegated_management_subnet_id").(string)

	id, err := parse.CassandraDatacenterID(d.Id())

	if err != nil {
		return fmt.Errorf("updating Cassandra MI Cluster %q (Datacenter: %q) - %+v", id.ClusterName, id.DatacenterName, err)
	}
	body := documentdb.DataCenterResource{
		//Location: &location,
		Properties: &documentdb.DataCenterResourceProperties{
			DelegatedSubnetID:  &delegatedSubnetId,
			NodeCount:          utils.Int32(int32(d.Get("node_count").(int))),
			DataCenterLocation: &location,
		},
	}

	future, err := client.CreateUpdate(ctx, id.ResourceGroup, id.ClusterName, id.DatacenterName, body)
	if err != nil {
		return fmt.Errorf("issuing create request for Cassandra MI  %q (Datacenter: %q): %+v", id.ClusterName, id.DatacenterName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on create/update future for Cassandra MI  %q (Datacenter: %q): %+v", id.ClusterName, id.DatacenterName, err)
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ClusterName, id.DatacenterName)

	if err != nil {
		return fmt.Errorf("making get request for Cassandra MI  %q (Datacenter: %q): %+v", id.ClusterName, id.DatacenterName, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("getting ID from Cassandra MI  %q (Datacenter: %q)", id.ClusterName, id.DatacenterName)
	}

	d.SetId(*resp.ID)

	return resourceCassandraMIDatacenterRead(d, meta)
}

func resourceCassandraMIDatacenterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CassandraDatacentersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	//datacenterName := d.Get("datacenter_name").(string)
	//id, err := parse.CassandraDatacenterID(d.Id() + "/dataCenters/" + datacenterName)
	id, err := parse.CassandraDatacenterID(d.Id())
	if err != nil {
		return err
	}
	resp, err := client.Get(ctx, id.ResourceGroup, id.ClusterName, id.DatacenterName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading Cassandra MI %q (Datacenter: %q) - removing from state", id.ClusterName, id.ClusterName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading Cassandra MI  %q (Datacenter: %q): %+v", id.ClusterName, id.ClusterName, err)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("cluster_name", id.ClusterName)
	d.Set("datacenter_name", id.DatacenterName)
	if props := resp.Properties; props != nil {
		if res := props; res != nil {
			d.Set("delegated_management_subnet_id", props.DelegatedSubnetID)
			d.Set("location", location.NormalizeNilable(props.DataCenterLocation))
			// nodeCountString := fmt.Sprint(*props.NodeCount)
			d.Set("node_count", int(*props.NodeCount))
		}
	}
	return nil
}

func resourceCassandraMIDatacenterDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CassandraDatacentersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	resourceGroup := d.Get("resource_group_name").(string)
	clusterName := d.Get("cluster_name").(string)
	datacenterName := d.Get("datacenter_name").(string)
	defer cancel()

	// id, err := parse.CassandraDatacenterID(d.Id() + "/dataCenters/" + datacenterName)
	// log.Println("************* id: " + d.Id() + "/dataCenters/" + datacenterName)
	// if err != nil {
	// 	return err
	// }

	future, err := client.Delete(ctx, resourceGroup, clusterName, datacenterName)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("deleting Cassandra MI  %q (Datacenter: %q): %+v", resourceGroup, clusterName, err)
		}
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("waiting on delete future for Cassandra MI  %q (Datacenter: %q): %+v", resourceGroup, clusterName, err)
	}

	return nil
}
