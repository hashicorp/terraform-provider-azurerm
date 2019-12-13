package azurerm

import (
	"fmt"
	"log"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"

	"github.com/Azure/azure-sdk-for-go/profiles/preview/preview/appplatform/mgmt/appplatform"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	azappplatform "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/appplatform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSpringCloudApp() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSpringCloudAppCreate,
		Read:   resourceArmSpringCloudAppRead,
		Update: resourceArmSpringCloudAppUpdate,
		Delete: resourceArmSpringCloudAppDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azappplatform.ValidateSpringCloudName,
			},

			"spring_cloud_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azappplatform.ValidateSpringCloudName,
			},

			"persistent_disk": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mount_path": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "/persistent",
							ValidateFunc: azappplatform.ValidateMountPath,
						},
						"size_in_gb": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntInSlice([]int{0, 50}),
						},
					},
				},
			},

			"public": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"temporary_disk": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mount_path": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "/tmp",
							ValidateFunc: azappplatform.ValidateMountPath,
						},
						"size_in_gb": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      5,
							ValidateFunc: validation.IntBetween(0, 5),
						},
					},
				},
			},

			"created_time": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmSpringCloudAppCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).AppPlatform.AppsClient
	ctx, cancel := timeouts.ForCreate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	springCloudName := d.Get("spring_cloud_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, springCloudName, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for present of existing Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q): %+v", name, springCloudName, resourceGroup, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_spring_cloud_app", *existing.ID)
		}
	}

	persistentDisk := d.Get("persistent_disk").([]interface{})
	public := d.Get("public").(bool)
	temporaryDisk := d.Get("temporary_disk").([]interface{})

	appResource := appplatform.AppResource{
		Properties: &appplatform.AppResourceProperties{
			PersistentDisk: expandArmSpringCloudAppPersistentDisk(persistentDisk),
			Public:         utils.Bool(public),
			TemporaryDisk:  expandArmSpringCloudAppTemporaryDisk(temporaryDisk),
		},
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, springCloudName, name, &appResource)
	if err != nil {
		return fmt.Errorf("Error creating Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q): %+v", name, springCloudName, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q): %+v", name, springCloudName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, springCloudName, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q): %+v", name, springCloudName, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q) ID", name, springCloudName, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmSpringCloudAppRead(d, meta)
}

func resourceArmSpringCloudAppRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).AppPlatform.AppsClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	springCloudName := id.Path["Spring"]
	name := id.Path["apps"]

	resp, err := client.Get(ctx, resourceGroup, springCloudName, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Spring Cloud App %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q): %+v", name, springCloudName, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("spring_cloud_name", springCloudName)
	if appResourceProperties := resp.Properties; appResourceProperties != nil {
		d.Set("created_time", (appResourceProperties.CreatedTime).String())
		if err := d.Set("persistent_disk", flattenArmSpringCloudAppPersistentDisk(appResourceProperties.PersistentDisk)); err != nil {
			return fmt.Errorf("Error setting `persistent_disk`: %+v", err)
		}
		if err := d.Set("temporary_disk", flattenArmSpringCloudAppTemporaryDisk(appResourceProperties.TemporaryDisk)); err != nil {
			return fmt.Errorf("Error setting `temporary_disk`: %+v", err)
		}
		d.Set("public", appResourceProperties.Public)
		d.Set("url", appResourceProperties.URL)
	}

	return nil
}

func resourceArmSpringCloudAppUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).AppPlatform.AppsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	springCloudName := d.Get("spring_cloud_name").(string)
	name := d.Get("name").(string)

	appResource, err := client.Get(ctx, resourceGroup, springCloudName, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(appResource.Response) {
			return fmt.Errorf("Error reading Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q): %+v", name, springCloudName, resourceGroup, err)
		}
		return fmt.Errorf("Error reading Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q): %+v", name, springCloudName, resourceGroup, err)
	}

	// if app does not have active deployment, it doesn't allow update, so just return
	if appResource.Properties.ActiveDeploymentName == nil || *appResource.Properties.ActiveDeploymentName == "" {
		return nil
	}

	persistentDisk := d.Get("persistent_disk").([]interface{})
	public := d.Get("public").(bool)
	temporaryDisk := d.Get("temporary_disk").([]interface{})

	appResource.Properties.PersistentDisk = expandArmSpringCloudAppPersistentDisk(persistentDisk)
	appResource.Properties.Public = utils.Bool(public)
	appResource.Properties.TemporaryDisk = expandArmSpringCloudAppTemporaryDisk(temporaryDisk)

	future, err := client.Update(ctx, resourceGroup, springCloudName, name, &appResource)
	if err != nil {
		return fmt.Errorf("Error updating Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q): %+v", name, springCloudName, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for update of Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q): %+v", name, springCloudName, resourceGroup, err)
	}

	return resourceArmSpringCloudAppRead(d, meta)
}

func resourceArmSpringCloudAppDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).AppPlatform.AppsClient
	ctx, cancel := timeouts.ForDelete(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	springCloudName := id.Path["Spring"]
	name := id.Path["apps"]

	if _, err := client.Delete(ctx, resourceGroup, springCloudName, name); err != nil {
		return fmt.Errorf("Error deleting Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q): %+v", name, springCloudName, resourceGroup, err)
	}

	return nil
}

func expandArmSpringCloudAppPersistentDisk(input []interface{}) *appplatform.PersistentDisk {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	sizeInGb := v["size_in_gb"].(int)
	mountPath := v["mount_path"].(string)

	result := appplatform.PersistentDisk{
		MountPath: utils.String(mountPath),
		SizeInGB:  utils.Int32(int32(sizeInGb)),
	}
	return &result
}

func expandArmSpringCloudAppTemporaryDisk(input []interface{}) *appplatform.TemporaryDisk {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	sizeInGb := v["size_in_gb"].(int)
	mountPath := v["mount_path"].(string)

	result := appplatform.TemporaryDisk{
		MountPath: utils.String(mountPath),
		SizeInGB:  utils.Int32(int32(sizeInGb)),
	}
	return &result
}

func flattenArmSpringCloudAppPersistentDisk(input *appplatform.PersistentDisk) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	if mountPath := input.MountPath; mountPath != nil {
		result["mount_path"] = *mountPath
	}
	if sizeInGb := input.SizeInGB; sizeInGb != nil {
		result["size_in_gb"] = int(*sizeInGb)
	}

	return []interface{}{result}
}

func flattenArmSpringCloudAppTemporaryDisk(input *appplatform.TemporaryDisk) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	if mountPath := input.MountPath; mountPath != nil {
		result["mount_path"] = *mountPath
	}
	if sizeInGb := input.SizeInGB; sizeInGb != nil {
		result["size_in_gb"] = int(*sizeInGb)
	}

	return []interface{}{result}
}
