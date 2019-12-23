package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/streamanalytics/mgmt/2016-03-01/streamanalytics"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmStreamAnalyticsOutputSql() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmStreamAnalyticsOutputSqlCreateUpdate,
		Read:   resourceArmStreamAnalyticsOutputSqlRead,
		Update: resourceArmStreamAnalyticsOutputSqlCreateUpdate,
		Delete: resourceArmStreamAnalyticsOutputSqlDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"stream_analytics_job_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"resource_group_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"server": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"database": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"table": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"user": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"password": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validate.NoEmptyStrings,
			},
		},
	}
}

func resourceArmStreamAnalyticsOutputSqlCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.OutputsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] Preparing arguments for Azure Stream Analytics SQL Output creation.")
	name := d.Get("name").(string)
	jobName := d.Get("stream_analytics_job_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, jobName, name)
		if err != nil && !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("Error checking for existing Azure Stream Analytics SQL Output %q (Job %q / Resource Group %q): %s", name, jobName, resourceGroup, err)
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_stream_analytics_output_mssql", *existing.ID)
		}
	}

	server := d.Get("server").(string)
	databaseName := d.Get("database").(string)
	tableName := d.Get("table").(string)
	sqlUser := d.Get("user").(string)
	sqlUserPassword := d.Get("password").(string)

	props := streamanalytics.Output{
		Name: utils.String(name),
		OutputProperties: &streamanalytics.OutputProperties{
			Datasource: &streamanalytics.AzureSQLDatabaseOutputDataSource{
				Type: streamanalytics.TypeMicrosoftSQLServerDatabase,
				AzureSQLDatabaseOutputDataSourceProperties: &streamanalytics.AzureSQLDatabaseOutputDataSourceProperties{
					Server:   utils.String(server),
					Database: utils.String(databaseName),
					User:     utils.String(sqlUser),
					Password: utils.String(sqlUserPassword),
					Table:    utils.String(tableName),
				},
			},
		},
	}

	if d.IsNewResource() {
		if _, err := client.CreateOrReplace(ctx, props, resourceGroup, jobName, name, "", ""); err != nil {
			return fmt.Errorf("Error Creating Stream Analytics Output SQL %q (Job %q / Resource Group %q): %+v", name, jobName, resourceGroup, err)
		}

		read, err := client.Get(ctx, resourceGroup, jobName, name)
		if err != nil {
			return fmt.Errorf("Error retrieving Stream Analytics Output SQL %q (Job %q / Resource Group %q): %+v", name, jobName, resourceGroup, err)
		}
		if read.ID == nil {
			return fmt.Errorf("Cannot read ID of Stream Analytics Output SQL %q (Job %q / Resource Group %q)", name, jobName, resourceGroup)
		}

		d.SetId(*read.ID)
	} else {
		if _, err := client.Update(ctx, props, resourceGroup, jobName, name, ""); err != nil {
			return fmt.Errorf("Error Updating Stream Analytics Output SQL %q (Job %q / Resource Group %q): %+v", name, jobName, resourceGroup, err)
		}
	}

	return resourceArmStreamAnalyticsOutputSqlRead(d, meta)
}

func resourceArmStreamAnalyticsOutputSqlRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.OutputsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	jobName := id.Path["streamingjobs"]
	name := id.Path["outputs"]

	resp, err := client.Get(ctx, resourceGroup, jobName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Output SQL %q was not found in Stream Analytics Job %q / Resource Group %q - removing from state!", name, jobName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Stream Output SQL %q (Stream Analytics Job %q / Resource Group %q): %+v", name, jobName, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("stream_analytics_job_name", jobName)

	if props := resp.OutputProperties; props != nil {
		v, ok := props.Datasource.AsAzureSQLDatabaseOutputDataSource()
		if !ok {
			return fmt.Errorf("Error converting Output Data Source to SQL Output: %+v", err)
		}

		d.Set("server", v.Server)
		d.Set("database", v.Database)
		d.Set("table", v.Table)
		d.Set("user", v.User)
	}

	return nil
}

func resourceArmStreamAnalyticsOutputSqlDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.OutputsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	jobName := id.Path["streamingjobs"]
	name := id.Path["outputs"]

	if resp, err := client.Delete(ctx, resourceGroup, jobName, name); err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("Error deleting Output SQL %q (Stream Analytics Job %q / Resource Group %q) %+v", name, jobName, resourceGroup, err)
		}
	}

	return nil
}
