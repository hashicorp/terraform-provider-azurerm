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
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCassandraMICluster() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCassandraMIClusterCreate,
		Read:   resourceCassandraMIClusterRead,
		Update: resourceCassandraMIClusterUpdate,
		Delete: resourceCassandraMIClusterDelete,

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

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CosmosAccountName,
			},

			"delegated_management_subnet_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CosmosEntityName,
			},

			"initial_cassandra_admin_password": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CosmosEntityName,
			},

			"throughput": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.CosmosThroughput,
			},

			"autoscale_settings": common.DatabaseAutoscaleSettingsSchema(),
		},
	}
}

func resourceCassandraMIClusterCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CassandraMIClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	//name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	clusterName := d.Get("cluster_name").(string)
	location := d.Get("location").(string)
	delegatedManagementSubnetId := d.Get("delegated_management_subnet_id").(string)
	initialCassandraAdminPassword := d.Get("initial_cassandra_admin_password").(string)

	existing, err := client.GetCluster(ctx, resourceGroup, clusterName, location)

	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of creating Cosmos Cassandra Keyspace %q (Account: %q): %+v", clusterName, location, err)
		}
	} else {
		if existing.ID == nil && *existing.ID == "" {
			return fmt.Errorf("generating import ID for Cosmos Cassandra Keyspace %q (Account: %q)", clusterName, location)
		}

		return tf.ImportAsExistsError("azurerm_cosmosdb_cassandra_keyspace", *existing.ID)
	}

	body := documentdb.CassandraMIClusterCreateUpdateParameters{
		Location: &location,
		CassandraMIClusterCreateUpdateProperties: &documentdb.CassandraMIClusterCreateUpdateProperties{
			DelegatedManagementSubnetId:   &delegatedManagementSubnetId,
			InitialCassandraAdminPassword: &initialCassandraAdminPassword,
		},
	}

	future, err := client.CreateUpdate(ctx, resourceGroup, clusterName, location, body)
	if err != nil {
		return fmt.Errorf("issuing create request for Cosmos Cassandra Keyspace %q (Account: %q): %+v", clusterName, clusterName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on create/update future for Cosmos Cassandra Keyspace %q (Account: %q): %+v", clusterName, clusterName, err)
	}

	resp, err := client.GetCluster(ctx, resourceGroup, clusterName, clusterName)

	if err != nil {
		return fmt.Errorf("making get request for Cosmos Cassandra Keyspace %q (Account: %q): %+v", clusterName, clusterName, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("getting ID from Cosmos Cassandra Keyspace %q (Account: %q)", clusterName, clusterName)
	}

	d.SetId(*resp.ID)

	return resourceCassandraMIClusterRead(d, meta)
}

func resourceCassandraMIClusterUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CassandraClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CassandraKeyspaceID(d.Id())
	if err != nil {
		return err
	}

	err = common.CheckForChangeFromAutoscaleAndManualThroughput(d)
	if err != nil {
		return fmt.Errorf("updating Cosmos Cassandra Keyspace %q (Account: %q) - %+v", id.Name, id.DatabaseAccountName, err)
	}

	db := documentdb.CassandraKeyspaceCreateUpdateParameters{
		CassandraKeyspaceCreateUpdateProperties: &documentdb.CassandraKeyspaceCreateUpdateProperties{
			Resource: &documentdb.CassandraKeyspaceResource{
				ID: &id.Name,
			},
			Options: &documentdb.CreateUpdateOptions{},
		},
	}

	future, err := client.CreateUpdateCassandraKeyspace(ctx, id.ResourceGroup, id.DatabaseAccountName, id.Name, db)
	if err != nil {
		return fmt.Errorf("issuing create/update request for Cosmos Cassandra Keyspace %q (Account: %q): %+v", id.ResourceGroup, id.DatabaseAccountName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on create/update future for Cosmos Cassandra Keyspace %q (Account: %q): %+v", id.ResourceGroup, id.DatabaseAccountName, err)
	}

	if common.HasThroughputChange(d) {
		throughputParameters := common.ExpandCosmosDBThroughputSettingsUpdateParameters(d)
		throughputFuture, err := client.UpdateCassandraKeyspaceThroughput(ctx, id.ResourceGroup, id.DatabaseAccountName, id.Name, *throughputParameters)
		if err != nil {
			if response.WasNotFound(throughputFuture.Response()) {
				return fmt.Errorf("setting Throughput for Cosmos Cassandra Keyspace %q (Account: %q): %+v - "+
					"If the collection has not been created with an initial throughput, you cannot configure it later.", id.Name, id.DatabaseAccountName, err)
			}
		}

		if err = throughputFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting on ThroughputUpdate future for Cosmos Cassandra Keyspace %q (Account: %q): %+v", id.Name, id.DatabaseAccountName, err)
		}
	}

	return resourceCosmosDbCassandraKeyspaceRead(d, meta)
}

func resourceCassandraMIClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CassandraMIClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CassandraClusterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetCluster(ctx, id.ResourceGroup, id.ClusterName, id.ClusterName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading Cosmos Cassandra Keyspace %q (Account: %q) - removing from state", id.ClusterName, id.ClusterName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading Cosmos Cassandra Keyspace %q (Account: %q): %+v", id.ClusterName, id.ClusterName, err)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("account_name", id.ClusterName)
	if props := resp.CassandraKeyspaceGetProperties; props != nil {
		if res := props.Resource; res != nil {
			d.Set("name", res.ID)
		}
	}
	return nil
}

func resourceCassandraMIClusterDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CassandraClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CassandraKeyspaceID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.DeleteCassandraKeyspace(ctx, id.ResourceGroup, id.DatabaseAccountName, id.Name)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("deleting Cosmos Cassandra Keyspace %q (Account: %q): %+v", id.Name, id.DatabaseAccountName, err)
		}
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("waiting on delete future for Cosmos Cassandra Keyspace %q (Account: %q): %+v", id.Name, id.DatabaseAccountName, err)
	}

	return nil
}
