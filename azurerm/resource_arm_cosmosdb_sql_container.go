package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2015-04-08/documentdb"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmCosmosDbSQLContainer() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmCosmosDbSQLContainerCreate,
		Read:   resourceArmCosmosDbSQLContainerRead,
		Delete: resourceArmCosmosDbSQLContainerDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CosmosEntityName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CosmosAccountName,
			},

			"database_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CosmosEntityName,
			},

			"partition_key_path": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"unique_key": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"paths": {
							Type:     schema.TypeSet,
							Required: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validate.NoEmptyStrings,
							},
						},
					},
				},
			},
		},
	}
}

func resourceArmCosmosDbSQLContainerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cosmos.DatabaseClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	database := d.Get("database_name").(string)
	account := d.Get("account_name").(string)
	partitionkeypaths := d.Get("partition_key_path").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.GetSQLContainer(ctx, resourceGroup, account, database, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of creating Cosmos SQL Container %s (Account: %s, Database:%s): %+v", name, account, database, err)
			}
		} else {
			id, err := azure.CosmosGetIDFromResponse(existing.Response)
			if err != nil {
				return fmt.Errorf("Error generating import ID for Cosmos SQL Container '%s' (Account: %s, Database:%s)", name, account, database)
			}

			return tf.ImportAsExistsError("azurerm_cosmosdb_sql_container", id)
		}
	}

	db := documentdb.SQLContainerCreateUpdateParameters{
		SQLContainerCreateUpdateProperties: &documentdb.SQLContainerCreateUpdateProperties{
			Resource: &documentdb.SQLContainerResource{
				ID: &name,
			},
			Options: map[string]*string{},
		},
	}

	if partitionkeypaths != "" {
		db.SQLContainerCreateUpdateProperties.Resource.PartitionKey = &documentdb.ContainerPartitionKey{
			Paths: &[]string{partitionkeypaths},
			Kind:  documentdb.PartitionKindHash,
		}
	}

	if keys := expandCosmosSQLContainerUniqueKeys(d.Get("unique_key").(*schema.Set)); keys != nil {
		db.SQLContainerCreateUpdateProperties.Resource.UniqueKeyPolicy = &documentdb.UniqueKeyPolicy{
			UniqueKeys: keys,
		}
	}

	future, err := client.CreateUpdateSQLContainer(ctx, resourceGroup, account, database, name, db)
	if err != nil {
		return fmt.Errorf("Error issuing create/update request for Cosmos SQL Container %s (Account: %s, Database:%s): %+v", name, account, database, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting on create/update future for Cosmos SQL Container %s (Account: %s, Database:%s): %+v", name, account, database, err)
	}

	resp, err := client.GetSQLContainer(ctx, resourceGroup, account, database, name)
	if err != nil {
		return fmt.Errorf("Error making get request for Cosmos SQL Container %s (Account: %s, Database:%s): %+v", name, account, database, err)
	}

	id, err := azure.CosmosGetIDFromResponse(resp.Response)
	if err != nil {
		return fmt.Errorf("Error retrieving the ID for Cosmos SQL Container '%s' (Account: %s, Database:%s) ID: %v", name, account, database, err)
	}
	d.SetId(id)

	return resourceArmCosmosDbSQLContainerRead(d, meta)
}

func resourceArmCosmosDbSQLContainerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cosmos.DatabaseClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseCosmosDatabaseContainerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetSQLContainer(ctx, id.ResourceGroup, id.Account, id.Database, id.Container)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading Cosmos SQL Container %s (Account %s) - removing from state", id.Database, id.Container)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Cosmos SQL Container %s (Account %s): %+v", id.Database, id.Container, err)
	}

	d.Set("name", id.Container)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("account_name", id.Account)
	d.Set("database_name", id.Database)

	if props := resp.SQLContainerProperties; props != nil {
		if pk := props.PartitionKey; pk != nil {
			if paths := pk.Paths; paths != nil {
				if len(*paths) > 1 {
					return fmt.Errorf("Error reading PartitionKey Paths, more then 1 returned")
				} else if len(*paths) == 1 {
					d.Set("partition_key_path", (*paths)[0])
				}
			}
		}

		if ukp := props.UniqueKeyPolicy; ukp != nil {
			if err := d.Set("unique_key", flattenCosmosSQLContainerUniqueKeys(ukp.UniqueKeys)); err != nil {
				return fmt.Errorf("Error setting `unique_key`: %+v", err)
			}
		}
	}

	return nil
}

func resourceArmCosmosDbSQLContainerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cosmos.DatabaseClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseCosmosDatabaseContainerID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.DeleteSQLContainer(ctx, id.ResourceGroup, id.Account, id.Database, id.Container)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error deleting Cosmos SQL Container %s (Account %s): %+v", id.Database, id.Container, err)
		}
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting on delete future for Cosmos SQL Container %s (Account %s): %+v", id.Database, id.Account, err)
	}

	return nil
}

func expandCosmosSQLContainerUniqueKeys(s *schema.Set) *[]documentdb.UniqueKey {
	i := s.List()
	if len(i) <= 0 || i[0] == nil {
		return nil
	}

	keys := make([]documentdb.UniqueKey, 0)
	for _, k := range i {
		key := k.(map[string]interface{})

		paths := key["paths"].(*schema.Set).List()
		if len(paths) == 0 {
			continue
		}

		keys = append(keys, documentdb.UniqueKey{
			Paths: utils.ExpandStringSlice(paths),
		})
	}

	return &keys
}

func flattenCosmosSQLContainerUniqueKeys(keys *[]documentdb.UniqueKey) *[]map[string]interface{} {
	if keys == nil {
		return nil
	}

	slice := make([]map[string]interface{}, 0)
	for _, k := range *keys {

		if k.Paths == nil {
			continue
		}

		slice = append(slice, map[string]interface{}{
			"paths": *k.Paths,
		})
	}

	return &slice
}
