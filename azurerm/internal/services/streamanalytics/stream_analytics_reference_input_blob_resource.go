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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/streamanalytics/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceStreamAnalyticsReferenceInputBlob() *schema.Resource {
	return &schema.Resource{
		Create: resourceStreamAnalyticsReferenceInputBlobCreate,
		Read:   resourceStreamAnalyticsReferenceInputBlobRead,
		Update: resourceStreamAnalyticsReferenceInputBlobUpdate,
		Delete: resourceStreamAnalyticsReferenceInputBlobDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.StreamInputID(id)
			return err
		}),

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
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
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

			"serialization": SchemaStreamAnalyticsStreamInputSerialization(),
		},
	}
}

func resourceStreamAnalyticsReferenceInputBlobCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.InputsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Stream Analytics Reference Input Blob creation.")
	resourceId := parse.NewStreamInputID(subscriptionId, d.Get("resource_group_name").(string), d.Get("stream_analytics_job_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceId.ResourceGroup, resourceId.StreamingjobName, resourceId.InputName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", resourceId, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_stream_analytics_reference_input_blob", resourceId.ID())
		}
	}

	props, err := getBlobReferenceInputProps(d)
	if err != nil {
		return fmt.Errorf("creating the input props for resource creation: %v", err)
	}

	if _, err := client.CreateOrReplace(ctx, props, resourceId.ResourceGroup, resourceId.StreamingjobName, resourceId.InputName, "", ""); err != nil {
		return fmt.Errorf("creating %s: %+v", resourceId, err)
	}

	d.SetId(resourceId.ID())
	return resourceStreamAnalyticsReferenceInputBlobRead(d, meta)
}

func resourceStreamAnalyticsReferenceInputBlobUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.InputsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Stream Analytics Reference Input Blob update.")
	id, err := parse.StreamInputID(d.Id())
	if err != nil {
		return err
	}

	props, err := getBlobReferenceInputProps(d)
	if err != nil {
		return fmt.Errorf("creating the input props for resource update: %v", err)
	}

	if _, err := client.Update(ctx, props, id.ResourceGroup, id.StreamingjobName, id.InputName, ""); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceStreamAnalyticsReferenceInputBlobRead(d, meta)
}

func resourceStreamAnalyticsReferenceInputBlobRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.InputsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StreamInputID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.StreamingjobName, id.InputName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.InputName)
	d.Set("stream_analytics_job_name", id.StreamingjobName)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := resp.Properties; props != nil {
		v, ok := props.AsReferenceInputProperties()
		if !ok {
			return fmt.Errorf("converting Reference Input Blob to a Reference Input: %+v", err)
		}

		blobInputDataSource, ok := v.Datasource.AsBlobReferenceInputDataSource()
		if !ok {
			return fmt.Errorf("converting Reference Input Blob to an Blob Stream Input: %+v", err)
		}

		d.Set("date_format", blobInputDataSource.DateFormat)
		d.Set("path_pattern", blobInputDataSource.PathPattern)
		d.Set("storage_container_name", blobInputDataSource.Container)
		d.Set("time_format", blobInputDataSource.TimeFormat)

		if accounts := blobInputDataSource.StorageAccounts; accounts != nil && len(*accounts) > 0 {
			account := (*accounts)[0]
			d.Set("storage_account_name", account.AccountName)
		}

		if err := d.Set("serialization", FlattenStreamAnalyticsStreamInputSerialization(v.Serialization)); err != nil {
			return fmt.Errorf("setting `serialization`: %+v", err)
		}
	}

	return nil
}

func resourceStreamAnalyticsReferenceInputBlobDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.InputsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StreamInputID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, id.ResourceGroup, id.StreamingjobName, id.InputName); err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("deleting %s: %+v", id, err)
		}
	}

	return nil
}

func getBlobReferenceInputProps(d *schema.ResourceData) (streamanalytics.Input, error) {
	name := d.Get("name").(string)
	containerName := d.Get("storage_container_name").(string)
	dateFormat := d.Get("date_format").(string)
	pathPattern := d.Get("path_pattern").(string)
	storageAccountKey := d.Get("storage_account_key").(string)
	storageAccountName := d.Get("storage_account_name").(string)
	timeFormat := d.Get("time_format").(string)

	serializationRaw := d.Get("serialization").([]interface{})
	serialization, err := ExpandStreamAnalyticsStreamInputSerialization(serializationRaw)
	if err != nil {
		return streamanalytics.Input{}, fmt.Errorf("expanding `serialization`: %+v", err)
	}

	props := streamanalytics.Input{
		Name: utils.String(name),
		Properties: &streamanalytics.ReferenceInputProperties{
			Type: streamanalytics.TypeReference,
			Datasource: &streamanalytics.BlobReferenceInputDataSource{
				Type: streamanalytics.TypeBasicReferenceInputDataSourceTypeMicrosoftStorageBlob,
				BlobReferenceInputDataSourceProperties: &streamanalytics.BlobReferenceInputDataSourceProperties{
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

	return props, nil
}
