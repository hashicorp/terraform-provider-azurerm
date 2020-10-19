package kusto

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-kusto-go/kusto"
	dataplaneKusto "github.com/Azure/azure-kusto-go/kusto"
	"github.com/Azure/azure-kusto-go/kusto/data/table"
	"github.com/Azure/azure-kusto-go/kusto/data/types"
	"github.com/Azure/azure-kusto-go/kusto/unsafe"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/kusto/parse"
	dataplaneTypes "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/kusto/types"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/kusto/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmKustoTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmKustoTableCreate,
		Read:   resourceArmKustoTableRead,
		Update: resourceArmKustoTableUpdate,
		Delete: resourceArmKustoTableDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.KustoTableID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.KustoTableName,
			},

			"database_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.KustoDatabaseID,
			},

			"column": {
				Type:     schema.TypeList,
				MinItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(types.Bool),
								string(types.DateTime),
								string(types.Decimal),
								string(types.Dynamic),
								string(types.GUID),
								string(types.Int),
								string(types.Long),
								string(types.Real),
								string(types.String),
								string(types.Timespan),
							}, false),
						},
					},
				},
			},

			"doc": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"folder": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceArmKustoTableCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.ClustersClient
	dbClient := meta.(*clients.Client).Kusto.DatabasesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Kusto Table creation.")

	name := d.Get("name").(string)
	databaseID := d.Get("database_id").(string)
	kustoDatabaseID, err := parse.KustoDatabaseID(databaseID)
	if err != nil {
		return err
	}

	cluster, err := client.Get(ctx, kustoDatabaseID.ResourceGroup, kustoDatabaseID.Cluster)
	if err != nil {
		return fmt.Errorf("retrieving Kusto Cluster %q (Resource Group %q): %+v", kustoDatabaseID.Cluster, kustoDatabaseID.ResourceGroup, err)
	}
	if cluster.ClusterProperties == nil || cluster.ClusterProperties.URI == nil {
		return fmt.Errorf("Kusto Cluster %q (Resource Group %q) URI property is nil or empty", kustoDatabaseID.Cluster, kustoDatabaseID.ResourceGroup)
	}

	_, err = dbClient.Get(ctx, kustoDatabaseID.ResourceGroup, kustoDatabaseID.Cluster, kustoDatabaseID.Name)
	if err != nil {
		return fmt.Errorf("retrieving Kusto Database %q (Resource Group %q, Cluster %q): %+v", kustoDatabaseID.Name, kustoDatabaseID.ResourceGroup, kustoDatabaseID.Cluster, err)
	}

	dataplaneClient, err := meta.(*clients.Client).Kusto.NewDataPlaneClient(*cluster.URI)
	if err != nil {
		return fmt.Errorf("init Kusto Data Plane Client: %+v", err)
	}

	id := fmt.Sprintf("%s/Tables/%s", databaseID, name)

	if d.IsNewResource() {
		rowIter, err := executeKustoMgmtStatement(ctx, dataplaneClient, kustoDatabaseID.Name, ".show tables")
		if err != nil {
			return fmt.Errorf("listing tables in Kusto Database %q (Resource Group %q, Cluster %q): %s", kustoDatabaseID.Name, kustoDatabaseID.ResourceGroup, kustoDatabaseID.Cluster, err)
		}
		defer rowIter.Stop()

		found := false
		err = rowIter.Do(func(row *table.Row) error {
			if len(row.Values) > 0 && row.Values[0].String() == name {
				found = true
				return nil
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("iterating tables in Kusto Database %q (Resource Group %q, Cluster %q): %s", kustoDatabaseID.Name, kustoDatabaseID.ResourceGroup, kustoDatabaseID.Cluster, err)
		}

		if found {
			return tf.ImportAsExistsError("azurerm_kusto_table", id)
		}
	}

	columns := d.Get("column").([]interface{})
	folder := d.Get("folder").(string)
	doc := d.Get("doc").(string)

	schema := createTableSchema(columns)

	createTableRawStmt := fmt.Sprintf(".create table %s (%s)", name, schema)
	if folder != "" || doc != "" {
		tableOptions := ""
		if doc != "" {
			tableOptions = "docstring=\"" + doc + "\""
		}
		if folder != "" {
			if tableOptions != "" {
				tableOptions += ","
			}
			tableOptions = tableOptions + "folder=\"" + folder + "\""
		}
		createTableRawStmt = fmt.Sprintf("%s with (%s)", createTableRawStmt, tableOptions)
	}

	if err := executeKustoMgmtStatementIgnoreResultSet(ctx, dataplaneClient, kustoDatabaseID.Name, createTableRawStmt); err != nil {
		return err
	}

	d.SetId(id)

	return resourceArmKustoTableRead(d, meta)
}

func resourceArmKustoTableUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.ClustersClient
	dbClient := meta.(*clients.Client).Kusto.DatabasesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Kusto Table update.")

	id, err := parse.KustoTableID(d.Id())
	if err != nil {
		return err
	}

	cluster, err := client.Get(ctx, id.ResourceGroup, id.Cluster)
	if err != nil {
		return fmt.Errorf("retrieving Kusto Cluster %q (Resource Group %q): %+v", id.Cluster, id.ResourceGroup, err)
	}
	if cluster.ClusterProperties == nil || cluster.ClusterProperties.URI == nil {
		return fmt.Errorf("Kusto Cluster %q (Resource Group %q) URI property is nil or empty", id.Cluster, id.ResourceGroup)
	}

	_, err = dbClient.Get(ctx, id.ResourceGroup, id.Cluster, id.Database)
	if err != nil {
		return fmt.Errorf("retrieving Kusto Database %q (Resource Group %q, Cluster %q): %+v", id.Database, id.ResourceGroup, id.Cluster, err)
	}

	dataplaneClient, err := meta.(*clients.Client).Kusto.NewDataPlaneClient(*cluster.URI)
	if err != nil {
		return fmt.Errorf("init Kusto Data Plane Client: %+v", err)
	}

	columns := d.Get("column").([]interface{})
	folder := d.Get("folder").(string)
	doc := d.Get("doc").(string)

	schema := createTableSchema(columns)

	alterTableStmt := fmt.Sprintf(".alter table %s (%s)", id.Name, schema)
	err = executeKustoMgmtStatementIgnoreResultSet(ctx, dataplaneClient, id.Database, alterTableStmt)
	if err != nil {
		return fmt.Errorf("updating Kusto Table %q (Resource Group %q, Cluster %q, Database %q): %+v", id.Name, id.ResourceGroup, id.Cluster, id.Database, err)
	}

	alterDocstringStmt := fmt.Sprintf(".alter table %s docstring \"%s\"", id.Name, doc)
	err = executeKustoMgmtStatementIgnoreResultSet(ctx, dataplaneClient, id.Database, alterDocstringStmt)
	if err != nil {
		return fmt.Errorf("updating Kusto Table %q (Resource Group %q, Cluster %q, Database %q): %+v", id.Name, id.ResourceGroup, id.Cluster, id.Database, err)
	}

	alterFolderStmt := fmt.Sprintf(".alter table  %s folder \"%s\"", id.Name, folder)
	err = executeKustoMgmtStatementIgnoreResultSet(ctx, dataplaneClient, id.Database, alterFolderStmt)
	if err != nil {
		return fmt.Errorf("updating Kusto Table %q (Resource Group %q, Cluster %q, Database %q): %+v", id.Name, id.ResourceGroup, id.Cluster, id.Database, err)
	}

	return resourceArmKustoTableRead(d, meta)
}

func resourceArmKustoTableRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.ClustersClient
	dbClient := meta.(*clients.Client).Kusto.DatabasesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.KustoTableID(d.Id())
	if err != nil {
		return err
	}

	cluster, err := client.Get(ctx, id.ResourceGroup, id.Cluster)
	if err != nil {
		if utils.ResponseWasNotFound(cluster.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Kusto Cluster %q (Resource Group %q): %+v", id.Cluster, id.ResourceGroup, err)
	}
	database, err := dbClient.Get(ctx, id.ResourceGroup, id.Cluster, id.Database)
	if err != nil {
		if utils.ResponseWasNotFound(database.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Kusto Database %q (Resource Group %q, Cluster %q): %+v", id.Database, id.ResourceGroup, id.Cluster, err)
	}

	if cluster.ClusterProperties == nil || cluster.ClusterProperties.URI == nil {
		return fmt.Errorf("Kusto Cluster %q (Resource Group %q) URI property is nil or empty", id.Cluster, id.ResourceGroup)
	}

	dataplaneClient, err := meta.(*clients.Client).Kusto.NewDataPlaneClient(*cluster.URI)
	if err != nil {
		return fmt.Errorf("init Kusto Data Plane Client: %+v", err)
	}

	showTablesStmt := fmt.Sprintf(".show tables | where TableName == \"%s\"", id.Name)
	iter, err := executeKustoMgmtStatement(ctx, dataplaneClient, id.Database, showTablesStmt)
	if err != nil {
		return fmt.Errorf("listing tables in Kusto Database %q (Resource Group %q, Cluster %q): %s", id.Database, id.ResourceGroup, id.Cluster, err)
	}
	defer iter.Stop()

	recs := []dataplaneTypes.KustoTableRecord{}
	_ = iter.Do(
		func(row *table.Row) error {
			rec := dataplaneTypes.KustoTableRecord{}
			if err := row.ToStruct(&rec); err != nil {
				return err
			}
			recs = append(recs, rec)
			return nil
		},
	)

	if len(recs) == 0 {
		d.SetId("")
		return nil
	}

	getSchemaStmt := fmt.Sprintf("%s | getschema", id.Name)
	iter, err = executeKustoQueryStatement(ctx, dataplaneClient, id.Database, getSchemaStmt)
	if err != nil {
		return fmt.Errorf("querying Kusto Table schema: %+v", err)
	}
	defer iter.Stop()

	columns := []dataplaneTypes.KustoTableColumnSchemaRecord{}
	err = iter.Do(
		func(row *table.Row) error {
			column := dataplaneTypes.KustoTableColumnSchemaRecord{}
			if err := row.ToStruct(&column); err != nil {
				return err
			}
			columns = append(columns, column)
			return nil
		},
	)

	if err != nil {
		return fmt.Errorf("processing Kusto Table schema: %+v", err)
	}

	databaseID := fmt.Sprintf("%s/Databases/%s", *cluster.ID, recs[0].DatabaseName)
	d.Set("database_id", databaseID)
	d.Set("name", recs[0].TableName)
	d.Set("folder", recs[0].Folder)
	d.Set("doc", recs[0].DocString)
	d.Set("column", flattenKustoTableColumns(columns))

	return nil
}

func resourceArmKustoTableDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.ClustersClient
	dbClient := meta.(*clients.Client).Kusto.DatabasesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.KustoTableID(d.Id())
	if err != nil {
		return err
	}

	cluster, err := client.Get(ctx, id.ResourceGroup, id.Cluster)
	if err != nil {
		return fmt.Errorf("retrieving Kusto Cluster %q (Resource Group %q): %+v", id.Cluster, id.ResourceGroup, err)
	}
	if cluster.ClusterProperties == nil || cluster.ClusterProperties.URI == nil {
		return fmt.Errorf("Kusto Cluster %q (Resource Group %q) URI property is nil or empty", id.Cluster, id.ResourceGroup)
	}

	_, err = dbClient.Get(ctx, id.ResourceGroup, id.Cluster, id.Database)
	if err != nil {
		return fmt.Errorf("retrieving Kusto Database %q (Resource Group %q, Cluster %q): %+v", id.Database, id.ResourceGroup, id.Cluster, err)
	}

	dataplaneClient, err := meta.(*clients.Client).Kusto.NewDataPlaneClient(*cluster.URI)
	if err != nil {
		return fmt.Errorf("init Kusto Data Plane Client: %+v", err)
	}

	dropTableStmt := fmt.Sprintf(".drop table %s ifexists", id.Name)
	err = executeKustoMgmtStatementIgnoreResultSet(ctx, dataplaneClient, id.Database, dropTableStmt)
	if err != nil {
		return fmt.Errorf("deleting Kusto Table %q (Cluster %q, Database %q): %+v", id.Name, id.Cluster, id.Database, err)
	}

	return nil
}

func flattenKustoTableColumns(columns []dataplaneTypes.KustoTableColumnSchemaRecord) []interface{} {
	if len(columns) == 0 {
		return nil
	}

	output := make([]interface{}, 0)
	for _, v := range columns {
		column := make(map[string]interface{})
		column["name"] = v.ColumnName
		column["type"] = v.ColumnType
		output = append(output, column)
	}

	return output
}

func executeKustoMgmtStatement(ctx context.Context, dataplaneClient *dataplaneKusto.Client, database string, rawStmt string) (*dataplaneKusto.RowIterator, error) {
	stmt := kusto.NewStmt("", kusto.UnsafeStmt(unsafe.Stmt{Add: true})).UnsafeAdd(rawStmt)
	rowIter, err := dataplaneClient.Mgmt(ctx, database, stmt)
	if err != nil {
		return nil, err
	}
	return rowIter, nil
}

func executeKustoQueryStatement(ctx context.Context, dataplaneClient *dataplaneKusto.Client, database string, rawStmt string) (*dataplaneKusto.RowIterator, error) {
	stmt := kusto.NewStmt("", kusto.UnsafeStmt(unsafe.Stmt{Add: true})).UnsafeAdd(rawStmt)
	rowIter, err := dataplaneClient.Query(ctx, database, stmt)
	if err != nil {
		return nil, err
	}
	return rowIter, nil
}

func executeKustoMgmtStatementIgnoreResultSet(ctx context.Context, dataplaneClient *dataplaneKusto.Client, database string, rawStmt string) error {
	rowIter, err := executeKustoMgmtStatement(ctx, dataplaneClient, database, rawStmt)
	if err != nil {
		return fmt.Errorf("Error querying Kusto: %+v", err)
	}
	defer rowIter.Stop()
	_ = rowIter.Do(
		func(row *table.Row) error {
			return nil
		},
	)
	return nil
}

func createTableSchema(columns []interface{}) string {
	schema := ""
	for i, v := range columns {
		column := v.(map[string]interface{})
		schema += fmt.Sprintf("%s:%s", column["name"], column["type"])
		if i < len(columns)-1 {
			schema += ","
		}
	}
	return schema
}
