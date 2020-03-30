package importexport

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/importexport/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceStorageImportExportJob() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmStorageImportExportJobRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.ImportExportJobName,
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"storage_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"export_blob_paths": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"export_blob_path_prefixes": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"export_blob_list_path": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"backup_drive_manifest": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"diagnostics_path": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"log_level": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"job_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"drive_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bit_locker_key": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"drive_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"manifest_file": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"manifest_hash": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"return_address": {
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
						"email": {
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
						"street_address1": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state_or_province": {
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

			"return_shipping": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"carrier_account_number": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"carrier_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"shipping_information": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"city": {
							Type:     schema.TypeString,
							Required: true,
						},
						"country_or_region": {
							Type:     schema.TypeString,
							Required: true,
						},
						"postal_code": {
							Type:     schema.TypeString,
							Required: true,
						},
						"recipient_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"state_or_province": {
							Type:     schema.TypeString,
							Required: true,
						},
						"street_address1": {
							Type:     schema.TypeString,
							Required: true,
						},
						"phone": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"street_address2": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceArmStorageImportExportJobRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ImportExport.JobClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, name, resourceGroup)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Azure Import/Export Job %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Azure Import/Export Job %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	d.SetId(*resp.ID)
	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if props := resp.Properties; props != nil {
		d.Set("backup_drive_manifest", props.BackupDriveManifest)
		d.Set("diagnostics_path", props.DiagnosticsPath)
		d.Set("job_type", props.JobType)
		d.Set("log_level", props.LogLevel)
		d.Set("storage_account_id", props.StorageAccountID)

		if export := props.Export; export != nil {
			d.Set("export_blob_list_path", export.BlobListblobPath)

			if exportBlobList := export.ExportBlobList; exportBlobList != nil {
				d.Set("export_blob_paths", utils.FlattenStringSlice(exportBlobList.BlobPath))
				d.Set("export_blob_path_prefixes", utils.FlattenStringSlice(exportBlobList.BlobPathPrefix))
			}
		}

		if err := d.Set("drive_info", flattenArmJobDriveInfo(props.DriveList, d)); err != nil {
			return fmt.Errorf("failure setting drive_info: %+v", err)
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
