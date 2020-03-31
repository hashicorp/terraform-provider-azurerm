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

func resourceArmStorageImportJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmStorageImportJobCreate,
		Read:   resourceArmStorageImportJobRead,
		Update: resourceArmStorageImportJobUpdate,
		Delete: resourceArmStorageImportJobDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImportThen(func(id string) error {
			_, err := parse.StorageImportExportJobID(id)
			return err
		}, importAzureImportExportJob(ImportJobType, "azurerm_import_job")),

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

			"drives": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 10,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"drive_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"bit_locker_key": {
							Type:         schema.TypeString,
							Sensitive:    true,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"manifest_file": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"manifest_hash": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
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
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Blue Dart",
								"DHL",
								"FedEx",
								"TNT",
								"UPS",
							}, false),
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
							ValidateFunc: validate.ImportExportJobEmail,
						},
						"phone": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.ImportExportJobPhone,
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
func resourceArmStorageImportJobCreate(d *schema.ResourceData, meta interface{}) error {
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
		return tf.ImportAsExistsError("azurerm_import_job", *existing.ID)
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	storageAccountId := d.Get("storage_account_id").(string)
	backupDriveManifest := d.Get("backup_drive_manifest").(bool)
	diagnosticsPath := d.Get("diagnostics_path").(string)
	logLevel := d.Get("log_level").(string)

	returnAddress := expandArmJobReturnAddress(d.Get("return_address").([]interface{}))
	returnShipping := expandArmJobReturnShipping(d.Get("return_shipping").([]interface{}))
	drives := expandArmJobDrives(d.Get("drives").([]interface{}))

	body := storageimportexport.PutJobParameters{
		Location: utils.String(location),
		Properties: &storageimportexport.JobDetails{
			BackupDriveManifest: utils.Bool(backupDriveManifest),
			DiagnosticsPath:     utils.String(diagnosticsPath),
			DriveList:           drives,
			JobType:             utils.String(ImportJobType),
			LogLevel:            utils.String(logLevel),
			ReturnAddress:       returnAddress,
			ReturnShipping:      returnShipping,
			StorageAccountID:    utils.String(storageAccountId),
		},
	}

	if _, err := client.Create(ctx, name, resourceGroup, body, ""); err != nil {
		return fmt.Errorf("failure creating Azure Import Job %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, name, resourceGroup)
	if err != nil {
		return fmt.Errorf("failure retrieving Azure Import Job %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("cannot read Azure Import Job %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)
	return resourceArmStorageImportJobRead(d, meta)
}

func resourceArmStorageImportJobUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ImportExport.JobClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	backupDriveManifest := d.Get("backup_drive_manifest").(bool)
	logLevel := d.Get("log_level").(string)

	returnAddress := expandArmJobReturnAddress(d.Get("return_address").([]interface{}))
	returnShipping := expandArmJobReturnShipping(d.Get("return_shipping").([]interface{}))
	drives := expandArmJobDrives(d.Get("drives").([]interface{}))

	body := storageimportexport.UpdateJobParameters{
		UpdateJobParametersProperties: &storageimportexport.UpdateJobParametersProperties{
			BackupDriveManifest: utils.Bool(backupDriveManifest),
			LogLevel:            utils.String(logLevel),
			DriveList:           drives,
			ReturnAddress:       returnAddress,
			ReturnShipping:      returnShipping,
		},
	}

	if _, err := client.Update(ctx, name, resourceGroup, body); err != nil {
		return fmt.Errorf("failure updating Azure Import Job %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, name, resourceGroup)
	if err != nil {
		return fmt.Errorf("failure retrieving Azure Import Job %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("Cannot read Azure Import Job %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)
	return resourceArmStorageImportJobRead(d, meta)
}

func resourceArmStorageImportJobRead(d *schema.ResourceData, meta interface{}) error {
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
			log.Printf("[INFO] Azure Import Job %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failure retrieving Azure Import Job %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	if props := resp.Properties; props != nil {
		d.Set("storage_account_id", props.StorageAccountID)
		d.Set("backup_drive_manifest", props.BackupDriveManifest)
		d.Set("diagnostics_path", props.DiagnosticsPath)
		d.Set("log_level", props.LogLevel)

		if err := d.Set("drives", flattenArmJobDrives(props.DriveList, d)); err != nil {
			return fmt.Errorf("failure setting drives: %+v", err)
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

func resourceArmStorageImportJobDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ImportExport.JobClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StorageImportExportJobID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.Name, id.ResourceGroup); err != nil {
		return fmt.Errorf("failure deleting Azure Import Job %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	return nil
}

func expandArmJobDrives(input []interface{}) *[]storageimportexport.DriveStatus {
	results := make([]storageimportexport.DriveStatus, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		driveId := v["drive_id"].(string)
		bitLockerKey := v["bit_locker_key"].(string)
		manifestFile := v["manifest_file"].(string)
		manifestHash := v["manifest_hash"].(string)
		result := storageimportexport.DriveStatus{
			DriveID:      utils.String(driveId),
			BitLockerKey: utils.String(bitLockerKey),
			ManifestFile: utils.String(manifestFile),
			ManifestHash: utils.String(manifestHash),
		}
		results = append(results, result)
	}
	return &results
}

func flattenArmJobDrives(input *[]storageimportexport.DriveStatus, d *schema.ResourceData) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	// prepare old state to find sensitive props not returned by API.
	oldDriveInfosRaw := d.Get("drives").([]interface{})

	for i, item := range *input {
		// prepare old state to find sensitive props not returned by API.
		oldDriveInfo := make(map[string]interface{})
		if len(oldDriveInfosRaw) > 0 {
			oldDriveInfo = oldDriveInfosRaw[i].(map[string]interface{})
		}

		// bit_locker_key returned by API are empty string
		// to avoid state diff, we get the props from old state
		var bitLockerKey string
		if v, ok := oldDriveInfo["bit_locker_key"]; ok {
			bitLockerKey = v.(string)
		}

		var driveId string
		if item.DriveID != nil {
			driveId = *item.DriveID
		}
		var manifestFile string
		if item.ManifestFile != nil {
			manifestFile = *item.ManifestFile
		}
		var manifestHash string
		if item.ManifestHash != nil {
			manifestHash = *item.ManifestHash
		}
		v := map[string]interface{}{
			"bit_locker_key": bitLockerKey,
			"drive_id":       driveId,
			"manifest_file":  manifestFile,
			"manifest_hash":  manifestHash,
		}
		results = append(results, v)
	}
	return results
}
