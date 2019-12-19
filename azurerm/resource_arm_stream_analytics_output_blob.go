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

func resourceArmStreamAnalyticsOutputBlob() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmStreamAnalyticsOutputBlobCreateUpdate,
		Read:   resourceArmStreamAnalyticsOutputBlobRead,
		Update: resourceArmStreamAnalyticsOutputBlobCreateUpdate,
		Delete: resourceArmStreamAnalyticsOutputBlobDelete,
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

			"resource_group_name": azure.SchemaResourceGroupName(),

			"date_format": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"path_pattern": {
				Type:     schema.TypeString,
				Required: true,
			},

			"storage_account_key": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"storage_account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"storage_container_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"time_format": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"serialization": azure.SchemaStreamAnalyticsOutputSerialization(),
		},
	}
}

func resourceArmStreamAnalyticsOutputBlobCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.OutputsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Stream Analytics Output Blob creation.")
	name := d.Get("name").(string)
	jobName := d.Get("stream_analytics_job_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, jobName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Stream Analytics Output Blob %q (Job %q / Resource Group %q): %s", name, jobName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_stream_analytics_output_blob", *existing.ID)
		}
	}

	containerName := d.Get("storage_container_name").(string)
	dateFormat := d.Get("date_format").(string)
	pathPattern := d.Get("path_pattern").(string)
	storageAccountKey := d.Get("storage_account_key").(string)
	storageAccountName := d.Get("storage_account_name").(string)
	timeFormat := d.Get("time_format").(string)

	serializationRaw := d.Get("serialization").([]interface{})
	serialization, err := azure.ExpandStreamAnalyticsOutputSerialization(serializationRaw)
	if err != nil {
		return fmt.Errorf("Error expanding `serialization`: %+v", err)
	}

	props := streamanalytics.Output{
		Name: utils.String(name),
		OutputProperties: &streamanalytics.OutputProperties{
			Datasource: &streamanalytics.BlobOutputDataSource{
				Type: streamanalytics.TypeMicrosoftStorageBlob,
				BlobOutputDataSourceProperties: &streamanalytics.BlobOutputDataSourceProperties{
					StorageAccounts: &[]streamanalytics.StorageAccount{
						{
							AccountKey:  utils.String(storageAccountKey),
							AccountName: utils.String(storageAccountName),
						},
					},
					Container:   utils.String(containerName),
					DateFormat:  utils.String(dateFormat),
					PathPattern: utils.String(pathPattern),
					TimeFormat:  utils.String(timeFormat),
				},
			},
			Serialization: serialization,
		},
	}

	if d.IsNewResource() {
		if _, err := client.CreateOrReplace(ctx, props, resourceGroup, jobName, name, "", ""); err != nil {
			return fmt.Errorf("Error Creating Stream Analytics Output Blob %q (Job %q / Resource Group %q): %+v", name, jobName, resourceGroup, err)
		}

		read, err := client.Get(ctx, resourceGroup, jobName, name)
		if err != nil {
			return fmt.Errorf("Error retrieving Stream Analytics Output Blob %q (Job %q / Resource Group %q): %+v", name, jobName, resourceGroup, err)
		}
		if read.ID == nil {
			return fmt.Errorf("Cannot read ID of Stream Analytics Output Blob %q (Job %q / Resource Group %q)", name, jobName, resourceGroup)
		}

		d.SetId(*read.ID)
	} else {
		if _, err := client.Update(ctx, props, resourceGroup, jobName, name, ""); err != nil {
			return fmt.Errorf("Error Updating Stream Analytics Output Blob %q (Job %q / Resource Group %q): %+v", name, jobName, resourceGroup, err)
		}
	}

	return resourceArmStreamAnalyticsOutputBlobRead(d, meta)
}

func resourceArmStreamAnalyticsOutputBlobRead(d *schema.ResourceData, meta interface{}) error {
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
			log.Printf("[DEBUG] Output Blob %q was not found in Stream Analytics Job %q / Resource Group %q - removing from state!", name, jobName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Stream Output EventHub %q (Stream Analytics Job %q / Resource Group %q): %+v", name, jobName, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("stream_analytics_job_name", jobName)

	if props := resp.OutputProperties; props != nil {
		v, ok := props.Datasource.AsBlobOutputDataSource()
		if !ok {
			return fmt.Errorf("Error converting Output Data Source to a Blob Output: %+v", err)
		}

		d.Set("date_format", v.DateFormat)
		d.Set("path_pattern", v.PathPattern)
		d.Set("storage_container_name", v.Container)
		d.Set("time_format", v.TimeFormat)

		if accounts := v.StorageAccounts; accounts != nil && len(*accounts) > 0 {
			account := (*accounts)[0]
			d.Set("storage_account_name", account.AccountName)
		}

		if err := d.Set("serialization", azure.FlattenStreamAnalyticsOutputSerialization(props.Serialization)); err != nil {
			return fmt.Errorf("Error setting `serialization`: %+v", err)
		}
	}

	return nil
}

func resourceArmStreamAnalyticsOutputBlobDelete(d *schema.ResourceData, meta interface{}) error {
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
			return fmt.Errorf("Error deleting Output Blob %q (Stream Analytics Job %q / Resource Group %q) %+v", name, jobName, resourceGroup, err)
		}
	}

	return nil
}
