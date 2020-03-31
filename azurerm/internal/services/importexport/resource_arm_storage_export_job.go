package importexport

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storageimportexport/mgmt/2016-11-01/storageimportexport"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/importexport/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/importexport/validate"
	storageValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmStorageExportJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmStorageExportJobCreate,
		Read:   resourceArmStorageExportJobRead,
		Update: resourceArmStorageExportJobUpdate,
		Delete: resourceArmStorageExportJobDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImportThen(func(id string) error {
			_, err := parse.StorageImportExportJobID(id)
			return err
		}, importAzureImportExportJob(ExportJobType, "azurerm_export_job")),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ImportExportJobName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"storage_account_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: storageValidate.StorageAccountID,
			},

			"export_blob_paths": {
				Type:         schema.TypeSet,
				Optional:     true,
				ForceNew:     true,
				AtLeastOneOf: []string{"export_blob_paths", "export_blob_path_prefixes", "export_blob_list_path"},
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},

			"export_blob_path_prefixes": {
				Type:         schema.TypeSet,
				Optional:     true,
				ForceNew:     true,
				AtLeastOneOf: []string{"export_blob_paths", "export_blob_path_prefixes", "export_blob_list_path"},
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},

			"export_blob_list_path": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				AtLeastOneOf: []string{"export_blob_paths", "export_blob_path_prefixes", "export_blob_list_path"},
			},

			"return_shipping": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"carrier_account_number": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"carrier_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"return_address": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"city": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"country_or_region": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"email": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"phone": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"postal_code": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"recipient_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"street_address1": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"street_address2": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"state_or_province": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"backup_drive_manifest": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"log_level": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Error",
				ValidateFunc: validation.StringInSlice([]string{
					"Error",
					"Verbose",
				}, false),
			},

			"diagnostics_path": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "waimportexport",
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"shipping_information": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"city": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"country_or_region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"phone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"postal_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"recipient_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state_or_province": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"street_address1": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"street_address2": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceArmStorageExportJobCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ImportExport.JobClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	existing, err := client.Get(ctx, name, resourceGroup)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("failure checking for present of existing Azure Import Export Job %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}
	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_export_job", *existing.ID)
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	storageAccountId := d.Get("storage_account_id").(string)

	blobPath := utils.ExpandStringSlice(d.Get("export_blob_paths").(*schema.Set).List())
	blobPathPrefix := utils.ExpandStringSlice(d.Get("export_blob_path_prefixes").(*schema.Set).List())
	blobListblobPath := d.Get("export_blob_list_path").(string)

	diagnosticsPath := d.Get("diagnostics_path").(string)
	logLevel := d.Get("log_level").(string)
	backupDriveManifest := d.Get("backup_drive_manifest").(bool)

	returnAddress := expandArmJobReturnAddress(d.Get("return_address").([]interface{}))
	returnShipping := expandArmJobReturnShipping(d.Get("return_shipping").([]interface{}))

	body := storageimportexport.PutJobParameters{
		Location: utils.String(location),
		Properties: &storageimportexport.JobDetails{
			BackupDriveManifest: utils.Bool(backupDriveManifest),
			DiagnosticsPath:     utils.String(diagnosticsPath),
			Export: &storageimportexport.Export{
				ExportBlobList: &storageimportexport.ExportBlobList{
					BlobPath:       blobPath,
					BlobPathPrefix: blobPathPrefix,
				},
				BlobListblobPath: utils.String(blobListblobPath),
			},
			JobType:          utils.String(ExportJobType),
			LogLevel:         utils.String(logLevel),
			StorageAccountID: utils.String(storageAccountId),
			ReturnAddress:    returnAddress,
			ReturnShipping:   returnShipping,
		},
	}

	if _, err := client.Create(ctx, name, resourceGroup, body, ""); err != nil {
		return fmt.Errorf("failure creating Azure Export Job %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, name, resourceGroup)
	if err != nil {
		return fmt.Errorf("failure retrieving Azure Export Job %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("cannot read Azure Export Job %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)
	return resourceArmStorageExportJobRead(d, meta)
}

func resourceArmStorageExportJobUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ImportExport.JobClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	body := storageimportexport.UpdateJobParameters{
		UpdateJobParametersProperties: &storageimportexport.UpdateJobParametersProperties{
			BackupDriveManifest: utils.Bool(d.Get("backup_drive_manifest").(bool)),
			LogLevel:            utils.String(d.Get("log_level").(string)),
			ReturnAddress:       expandArmJobReturnAddress(d.Get("return_address").([]interface{})),
			ReturnShipping:      expandArmJobReturnShipping(d.Get("return_shipping").([]interface{})),
		},
	}

	if _, err := client.Update(ctx, name, resourceGroup, body); err != nil {
		return fmt.Errorf("failure updating Azure Export Job %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, name, resourceGroup)
	if err != nil {
		return fmt.Errorf("failure retrieving Azure Export Job %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("Cannot read Azure Export Job %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)
	return resourceArmStorageExportJobRead(d, meta)
}

func resourceArmStorageExportJobRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ImportExport.JobClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StorageImportExportJobID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.Name, id.ResourceGroup)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Azure Import Export Job %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failure retrieving Azure Import Export Job %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	if props := resp.Properties; props != nil {
		d.Set("storage_account_id", props.StorageAccountID)
		d.Set("backup_drive_manifest", props.BackupDriveManifest)
		d.Set("diagnostics_path", props.DiagnosticsPath)
		d.Set("log_level", props.LogLevel)
		if export := props.Export; export != nil {
			d.Set("export_blob_list_path", export.BlobListblobPath)

			if exportBlobList := export.ExportBlobList; exportBlobList != nil {
				d.Set("export_blob_paths", utils.FlattenStringSlice(exportBlobList.BlobPath))
				d.Set("export_blob_path_prefixes", utils.FlattenStringSlice(exportBlobList.BlobPathPrefix))
			}
		}

		if err := d.Set("return_address", flattenArmJobReturnAddress(props.ReturnAddress)); err != nil {
			return fmt.Errorf("failure setting return_address: %+v", err)
		}
		if err := d.Set("return_shipping", flattenArmJobReturnShipping(props.ReturnShipping)); err != nil {
			return fmt.Errorf("failure setting return_shipping: %+v", err)
		}
		if err := d.Set("shipping_information", flattenArmJobShippingInformation(props.ShippingInformation)); err != nil {
			return fmt.Errorf("failure setting shipping_information: %+v", err)
		}
	}
	return nil
}

func resourceArmStorageExportJobDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ImportExport.JobClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StorageImportExportJobID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.Name, id.ResourceGroup); err != nil {
		return fmt.Errorf("failure deleting Azure Export Job %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	return nil
}
