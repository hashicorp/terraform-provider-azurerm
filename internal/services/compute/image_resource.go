// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/images"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceImage() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceImageCreateUpdate,
		Read:   resourceImageRead,
		Update: resourceImageCreateUpdate,
		Delete: resourceImageDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := images.ParseImageID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(90 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(90 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(90 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"zone_resilient": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"hyper_v_generation": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(images.HyperVGenerationTypesVOne),
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(images.HyperVGenerationTypesVOne),
					string(images.HyperVGenerationTypesVTwo),
				}, false),
			},

			"source_virtual_machine_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: commonids.ValidateVirtualMachineID,
			},

			"os_disk": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				ForceNew: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"os_type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(images.OperatingSystemTypesLinux),
								string(images.OperatingSystemTypesWindows),
							}, false),
						},

						"os_state": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(images.OperatingSystemStateTypesGeneralized),
								string(images.OperatingSystemStateTypesSpecialized),
							}, false),
						},

						"managed_disk_id": {
							Type:             pluginsdk.TypeString,
							Computed:         true,
							Optional:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc:     commonids.ValidateManagedDiskID,
						},

						"blob_uri": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
						},

						"caching": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(images.CachingTypesNone),
							ValidateFunc: validation.StringInSlice([]string{
								string(images.CachingTypesNone),
								string(images.CachingTypesReadOnly),
								string(images.CachingTypesReadWrite),
							}, false),
						},

						"size_gb": {
							Type:         pluginsdk.TypeInt,
							Computed:     true,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.NoZeroValues,
						},

						"disk_encryption_set_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validate.DiskEncryptionSetID,
						},

						"storage_type": {
							Type:         pluginsdk.TypeString,
							Description:  "The type of storage disk",
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice(images.PossibleValuesForStorageAccountTypes(), false),
						},
					},
				},
			},

			"data_disk": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"lun": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
						},

						"managed_disk_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: commonids.ValidateManagedDiskID,
						},

						"blob_uri": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
						},

						"caching": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(images.CachingTypesNone),
							ValidateFunc: validation.StringInSlice([]string{
								string(images.CachingTypesNone),
								string(images.CachingTypesReadOnly),
								string(images.CachingTypesReadWrite),
							}, false),
						},

						"size_gb": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.NoZeroValues,
						},

						"disk_encryption_set_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validate.DiskEncryptionSetID,
						},

						"storage_type": {
							Type:         pluginsdk.TypeString,
							Description:  "The type of storage disk",
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice(images.PossibleValuesForStorageAccountTypes(), false),
						},
					},
				},
			},

			"tags": commonschema.Tags(),
		},
	}

	if !features.FourPointOhBeta() {
		delete(resource.Schema["os_disk"].Elem.(*pluginsdk.Resource).Schema, "storage_type")
		delete(resource.Schema["data_disk"].Elem.(*pluginsdk.Resource).Schema, "storage_type")
	}

	return resource
}

func resourceImageCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.ImagesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := images.NewImageID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id, images.DefaultGetOperationOptions())
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_image", id.ID())
		}
	}

	props := images.ImageProperties{
		HyperVGeneration: pointer.To(images.HyperVGenerationTypes(d.Get("hyper_v_generation").(string))),
	}

	sourceVM := images.SubResource{}
	if v, ok := d.GetOk("source_virtual_machine_id"); ok {
		sourceVM.Id = pointer.To(v.(string))
	}

	storageProfile := images.ImageStorageProfile{
		OsDisk:        expandImageOSDisk(d.Get("os_disk").([]interface{})),
		DataDisks:     expandImageDataDisks(d.Get("data_disk").([]interface{})),
		ZoneResilient: utils.Bool(d.Get("zone_resilient").(bool)),
	}

	// either source VM or storage profile can be specified, but not both
	if sourceVM.Id == nil {
		// if both sourceVM and storageProfile are empty, return an error
		if storageProfile.OsDisk == nil && (storageProfile.DataDisks == nil || len(*storageProfile.DataDisks) == 0) {
			return fmt.Errorf("[ERROR] Cannot create image when both source VM and storage profile are empty")
		}

		props.StorageProfile = &storageProfile
	} else {
		// creating an image from source VM
		props.SourceVirtualMachine = &sourceVM
	}

	payload := images.Image{
		Location:   location.Normalize(d.Get("location").(string)),
		Properties: &props,
		Tags:       tags.Expand(d.Get("tags").(map[string]interface{})),
	}
	if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceImageRead(d, meta)
}

func resourceImageRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.ImagesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := images.ParseImageID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, images.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.ImageName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			hyperVGeneration := ""
			if props.HyperVGeneration != nil {
				hyperVGeneration = string(*props.HyperVGeneration)
			}
			d.Set("hyper_v_generation", hyperVGeneration)

			// either source VM or storage profile can be specified, but not both
			if props.SourceVirtualMachine != nil && props.SourceVirtualMachine.Id != nil {
				d.Set("source_virtual_machine_id", pointer.From(props.SourceVirtualMachine.Id))
			} else {
				if err := d.Set("os_disk", flattenImageOSDisk(props.StorageProfile)); err != nil {
					return fmt.Errorf("setting `os_disk`: %+v", err)
				}
				if err := d.Set("data_disk", flattenImageDataDisks(props.StorageProfile)); err != nil {
					return fmt.Errorf("setting `data_disk`: %+v", err)
				}
				zoneResilient := false
				if props.StorageProfile != nil && props.StorageProfile.ZoneResilient != nil {
					zoneResilient = *props.StorageProfile.ZoneResilient
				}
				d.Set("zone_resilient", zoneResilient)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	return nil
}

func resourceImageDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.ImagesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := images.ParseImageID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandImageOSDisk(input []interface{}) *images.ImageOSDisk {
	if len(input) > 0 {
		config := input[0].(map[string]interface{})

		out := &images.ImageOSDisk{}

		if v := config["os_type"].(string); v != "" {
			out.OsType = images.OperatingSystemTypes(v)
		}

		if v := config["os_state"].(string); v != "" {
			out.OsState = images.OperatingSystemStateTypes(v)
		}
		managedDiskID := config["managed_disk_id"].(string)
		if managedDiskID != "" {
			managedDisk := &images.SubResource{
				Id: &managedDiskID,
			}
			out.ManagedDisk = managedDisk
		}

		blobURI := config["blob_uri"].(string)
		out.BlobUri = &blobURI

		if v := config["caching"].(string); v != "" {
			out.Caching = pointer.To(images.CachingTypes(v))
		}

		if size := config["size_gb"]; size != 0 {
			out.DiskSizeGB = pointer.To(int64(size.(int)))
		}

		if id := config["disk_encryption_set_id"].(string); id != "" {
			out.DiskEncryptionSet = &images.SubResource{
				Id: pointer.To(id),
			}
		}

		if features.FourPointOhBeta() {
			out.StorageAccountType = pointer.To(images.StorageAccountTypes(config["storage_type"].(string)))
		}
		return out
	}

	return nil
}

func expandImageDataDisks(disks []interface{}) *[]images.ImageDataDisk {
	output := make([]images.ImageDataDisk, 0)
	for _, diskConfig := range disks {
		config := diskConfig.(map[string]interface{})

		item := images.ImageDataDisk{
			BlobUri: pointer.To(config["blob_uri"].(string)),
			Lun:     int64(config["lun"].(int)),
		}

		if size := config["size_gb"]; size != 0 {
			item.DiskSizeGB = pointer.To(int64(size.(int)))
		}

		if v := config["caching"].(string); v != "" {
			item.Caching = pointer.To(images.CachingTypes(v))
		}

		if managedDiskID := config["managed_disk_id"].(string); managedDiskID != "" {
			managedDisk := &images.SubResource{
				Id: &managedDiskID,
			}
			item.ManagedDisk = managedDisk
		}

		if id := config["disk_encryption_set_id"].(string); id != "" {
			item.DiskEncryptionSet = &images.SubResource{
				Id: pointer.To(id),
			}
		}

		if features.FourPointOhBeta() {
			item.StorageAccountType = pointer.To(images.StorageAccountTypes(config["storage_type"].(string)))
		}

		output = append(output, item)
	}

	return &output
}

func flattenImageOSDisk(input *images.ImageStorageProfile) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		if v := input.OsDisk; v != nil {
			blobUri := ""
			if uri := v.BlobUri; uri != nil {
				blobUri = *uri
			}
			caching := ""
			if v.Caching != nil {
				caching = string(*v.Caching)
			}
			diskSizeGB := 0
			if v.DiskSizeGB != nil {
				diskSizeGB = int(*v.DiskSizeGB)
			}
			managedDiskId := ""
			if disk := v.ManagedDisk; disk != nil && disk.Id != nil {
				managedDiskId = *disk.Id
			}
			diskEncryptionSetId := ""
			if set := v.DiskEncryptionSet; set != nil && set.Id != nil {
				encryptionId, _ := commonids.ParseDiskEncryptionSetIDInsensitively(*set.Id)
				diskEncryptionSetId = encryptionId.ID()
			}

			properties := map[string]interface{}{
				"blob_uri":               blobUri,
				"caching":                caching,
				"managed_disk_id":        managedDiskId,
				"os_type":                string(v.OsType),
				"os_state":               string(v.OsState),
				"size_gb":                diskSizeGB,
				"disk_encryption_set_id": diskEncryptionSetId,
			}

			if features.FourPointOhBeta() {
				storageType := ""
				if v.StorageAccountType != nil {
					storageType = string(*v.StorageAccountType)
				}
				properties["storage_type"] = storageType
			}

			output = append(output, properties)
		}
	}

	return output
}

func flattenImageDataDisks(input *images.ImageStorageProfile) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		if v := input.DataDisks; v != nil {
			for _, disk := range *input.DataDisks {
				blobUri := ""
				if disk.BlobUri != nil {
					blobUri = *disk.BlobUri
				}
				caching := ""
				if disk.Caching != nil {
					caching = string(*disk.Caching)
				}
				diskSizeGb := 0
				if disk.DiskSizeGB != nil {
					diskSizeGb = int(*disk.DiskSizeGB)
				}
				managedDiskId := ""
				if disk.ManagedDisk != nil && disk.ManagedDisk.Id != nil {
					managedDiskId = *disk.ManagedDisk.Id
				}
				diskEncryptionSetId := ""
				if set := disk.DiskEncryptionSet; set != nil && set.Id != nil {
					encryptionId, _ := commonids.ParseDiskEncryptionSetIDInsensitively(*set.Id)
					diskEncryptionSetId = encryptionId.ID()
				}

				properties := map[string]interface{}{
					"blob_uri":               blobUri,
					"caching":                caching,
					"lun":                    int(disk.Lun),
					"managed_disk_id":        managedDiskId,
					"size_gb":                diskSizeGb,
					"disk_encryption_set_id": diskEncryptionSetId,
				}

				if features.FourPointOhBeta() {
					storageType := ""
					if disk.StorageAccountType != nil {
						storageType = string(*disk.StorageAccountType)
					}
					properties["storage_type"] = storageType
				}

				output = append(output, properties)
			}
		}
	}

	return output
}
