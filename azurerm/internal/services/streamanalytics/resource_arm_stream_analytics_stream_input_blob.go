package streamanalytics

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/streamanalytics/mgmt/2016-03-01/streamanalytics"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmStreamAnalyticsStreamInputBlob() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmStreamAnalyticsStreamInputBlobCreateUpdate,
		Read:   resourceArmStreamAnalyticsStreamInputBlobRead,
		Update: resourceArmStreamAnalyticsStreamInputBlobCreateUpdate,
		Delete: resourceArmStreamAnalyticsStreamInputBlobDelete,
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
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"stream_analytics_job_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"date_format": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"path_pattern": {
				Type:     schema.TypeString,
				Required: true,
			},

			"storage_account_key": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"storage_account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"storage_container_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"time_format": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"serialization": azure.SchemaStreamAnalyticsStreamInputSerialization(),
		},
	}
}

func resourceArmStreamAnalyticsStreamInputBlobCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.InputsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Stream Analytics Stream Input Blob creation.")
	name := d.Get("name").(string)
	jobName := d.Get("stream_analytics_job_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, jobName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Stream Analytics Stream Input %q (Job %q / Resource Group %q): %s", name, jobName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_stream_analytics_stream_input_blob", *existing.ID)
		}
	}

	containerName := d.Get("storage_container_name").(string)
	dateFormat := d.Get("date_format").(string)
	pathPattern := d.Get("path_pattern").(string)
	storageAccountKey := d.Get("storage_account_key").(string)
	storageAccountName := d.Get("storage_account_name").(string)
	timeFormat := d.Get("time_format").(string)

	serializationRaw := d.Get("serialization").([]interface{})
	serialization, err := azure.ExpandStreamAnalyticsStreamInputSerialization(serializationRaw)
	if err != nil {
		return fmt.Errorf("Error expanding `serialization`: %+v", err)
	}

	props := streamanalytics.Input{
		Name: utils.String(name),
		Properties: &streamanalytics.StreamInputProperties{
			Type: streamanalytics.TypeStream,
			Datasource: &streamanalytics.BlobStreamInputDataSource{
				Type: streamanalytics.TypeBasicStreamInputDataSourceTypeMicrosoftStorageBlob,
				BlobStreamInputDataSourceProperties: &streamanalytics.BlobStreamInputDataSourceProperties{
					Container:   utils.String(containerName),
					DateFormat:  utils.String(dateFormat),
					PathPattern: utils.String(pathPattern),
					TimeFormat:  utils.String(timeFormat),
					StorageAccounts: &[]streamanalytics.StorageAccount{
						{
							AccountName: utils.String(storageAccountName),
							AccountKey:  utils.String(storageAccountKey),
						},
					},
				},
			},
			Serialization: serialization,
		},
	}

	if d.IsNewResource() {
		if _, err := client.CreateOrReplace(ctx, props, resourceGroup, jobName, name, "", ""); err != nil {
			return fmt.Errorf("Error Creating Stream Analytics Stream Input Blob %q (Job %q / Resource Group %q): %+v", name, jobName, resourceGroup, err)
		}

		read, err := client.Get(ctx, resourceGroup, jobName, name)
		if err != nil {
			return fmt.Errorf("Error retrieving Stream Analytics Stream Input Blob %q (Job %q / Resource Group %q): %+v", name, jobName, resourceGroup, err)
		}
		if read.ID == nil {
			return fmt.Errorf("Cannot read ID of Stream Analytics Stream Input Blob %q (Job %q / Resource Group %q)", name, jobName, resourceGroup)
		}

		d.SetId(*read.ID)
	} else if _, err := client.Update(ctx, props, resourceGroup, jobName, name, ""); err != nil {
		return fmt.Errorf("Error Updating Stream Analytics Stream Input Blob %q (Job %q / Resource Group %q): %+v", name, jobName, resourceGroup, err)
	}

	return resourceArmStreamAnalyticsStreamInputBlobRead(d, meta)
}

func resourceArmStreamAnalyticsStreamInputBlobRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.InputsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	jobName := id.Path["streamingjobs"]
	name := id.Path["inputs"]

	resp, err := client.Get(ctx, resourceGroup, jobName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Stream Input Blob %q was not found in Stream Analytics Job %q / Resource Group %q - removing from state!", name, jobName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Stream Input Blob %q (Stream Analytics Job %q / Resource Group %q): %+v", name, jobName, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("stream_analytics_job_name", jobName)

	if props := resp.Properties; props != nil {
		v, ok := props.AsStreamInputProperties()
		if !ok {
			return fmt.Errorf("Error converting Stream Input Blob to an Stream Input: %+v", err)
		}

		eventHub, ok := v.Datasource.AsBlobStreamInputDataSource()
		if !ok {
			return fmt.Errorf("Error converting Stream Input Blob to an Blob Stream Input: %+v", err)
		}

		d.Set("date_format", eventHub.DateFormat)
		d.Set("path_pattern", eventHub.PathPattern)
		d.Set("storage_container_name", eventHub.Container)
		d.Set("time_format", eventHub.TimeFormat)

		if accounts := eventHub.StorageAccounts; accounts != nil && len(*accounts) > 0 {
			account := (*accounts)[0]
			d.Set("storage_account_name", account.AccountName)
		}

		if err := d.Set("serialization", azure.FlattenStreamAnalyticsStreamInputSerialization(v.Serialization)); err != nil {
			return fmt.Errorf("Error setting `serialization`: %+v", err)
		}
	}

	return nil
}

func resourceArmStreamAnalyticsStreamInputBlobDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.InputsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	jobName := id.Path["streamingjobs"]
	name := id.Path["inputs"]

	if resp, err := client.Delete(ctx, resourceGroup, jobName, name); err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("Error deleting Stream Input Blob %q (Stream Analytics Job %q / Resource Group %q) %+v", name, jobName, resourceGroup, err)
		}
	}

	return nil
}
