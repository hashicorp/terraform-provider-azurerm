package azurerm

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/preview/preview/appplatform/mgmt/appplatform"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	azappplatform "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/appplatform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"

	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/file/files"
)

func resourceArmSpringCloudDeployment() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSpringCloudDeploymentCreate,
		Read:   resourceArmSpringCloudDeploymentRead,
		Update: resourceArmSpringCloudDeploymentUpdate,
		Delete: resourceArmSpringCloudDeploymentDelete,

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

			"spring_cloud_app_name": {
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

			"cpu": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 4),
			},

			"memory_in_gb": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 8),
			},

			"instance_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 20),
			},

			"jvm_options": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"env": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"runtime_version": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(appplatform.Java8),
					string(appplatform.Java11),
				}, true),
				Default: string(appplatform.Java8),
			},

			"jar_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceArmSpringCloudDeploymentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).AppPlatform.DeploymentsClient
	ctx, cancel := timeouts.ForCreate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	appName := d.Get("spring_cloud_app_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	springCloudName := d.Get("spring_cloud_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, springCloudName, appName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for present of existing Spring Cloud Deployment %q (Spring Cloud Service %q / App Name %q / Resource Group %q): %+v", appName, springCloudName, appName, resourceGroup, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_spring_cloud_deployment", *existing.ID)
		}
	}

	cpu := int32(d.Get("cpu").(int))
	memoryInGB := int32(d.Get("memory_in_gb").(int))
	jvmOptions := d.Get("jvm_options").(string)
	instanceCount := int32(d.Get("instance_count").(int))
	runtimeVersion := d.Get("runtime_version").(string)

	userSourceInfo, err := expandSpringCloudDeploymentUserSourceInfo(d, meta, resourceGroup, springCloudName, appName)
	if err != nil {
		return err
	}

	appResource := appplatform.DeploymentResource{
		Properties: &appplatform.DeploymentResourceProperties{
			AppName: utils.String(appName),
			DeploymentSettings: &appplatform.DeploymentSettings{
				CPU:                  &cpu,
				MemoryInGB:           &memoryInGB,
				JvmOptions:           &jvmOptions,
				InstanceCount:        &instanceCount,
				EnvironmentVariables: expandSpringCloudDeploymentEnv(d),
				RuntimeVersion:       appplatform.RuntimeVersion(runtimeVersion),
			},
			Source: userSourceInfo,
		},
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, springCloudName, appName, name, &appResource)
	if err != nil {
		return fmt.Errorf("Error creating Spring Cloud Deployment %q (Spring Cloud Service %q / App Name %q / Resource Group %q): %+v", name, springCloudName, appName, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Spring Cloud Deployment %q (Spring Cloud Service %q / App Name %q /  Resource Group %q): %+v", name, springCloudName, appName, resourceGroup, err)
	}

	resp, err := future.Result(*client)
	if err != nil {
		return fmt.Errorf("Error retrieving Spring Cloud Deployment %q (Spring Cloud Service %q / App Name %q / Resource Group %q): %+v", name, springCloudName, appName, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot get Spring Cloud Deployment %q (Spring Cloud Service %q / App Name %q / Resource Group %q): %+v", name, springCloudName, appName, resourceGroup, err)
	}
	d.SetId(*resp.ID)

	return resourceArmSpringCloudDeploymentRead(d, meta)
}

func resourceArmSpringCloudDeploymentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).AppPlatform.DeploymentsClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	springCloudName := id.Path["Spring"]
	appName := id.Path["apps"]
	name := id.Path["deployments"]

	resp, err := client.Get(ctx, resourceGroup, springCloudName, appName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Spring Cloud Deployments %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Spring Cloud Deployment %q (Spring Cloud Service %q / App Name %q / Resource Group %q): %+v", name, springCloudName, appName, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("spring_cloud_name", springCloudName)
	d.Set("spring_cloud_app_name", appName)
	if deploymentSettings := resp.Properties.DeploymentSettings; deploymentSettings != nil {
		d.Set("cpu", deploymentSettings.CPU)
		d.Set("memory_in_gb", deploymentSettings.MemoryInGB)
		d.Set("jvm_options", deploymentSettings.JvmOptions)
		d.Set("instance_count", deploymentSettings.InstanceCount)
		d.Set("runtime_version", deploymentSettings.RuntimeVersion)
		if err := d.Set("env", flattenSpringCloudDeploymentEnv(deploymentSettings.EnvironmentVariables)); err != nil {
			return fmt.Errorf("Error setting `env`: %+v", err)
		}
	}

	return nil
}

func resourceArmSpringCloudDeploymentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).AppPlatform.DeploymentsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	appName := d.Get("spring_cloud_app_name").(string)
	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	springCloudName := d.Get("spring_cloud_name").(string)

	cpu := d.Get("cpu").(int32)
	memoryInGB := d.Get("memory_in_gb").(int32)
	jvmOptions := d.Get("jvm_options").(string)
	instanceCount := d.Get("instance_count").(int32)
	runtimeVersion := d.Get("runtime_version").(string)

	appResource := appplatform.DeploymentResource{
		Properties: &appplatform.DeploymentResourceProperties{
			AppName: utils.String(appName),
			DeploymentSettings: &appplatform.DeploymentSettings{
				CPU:                  &cpu,
				MemoryInGB:           &memoryInGB,
				JvmOptions:           &jvmOptions,
				InstanceCount:        &instanceCount,
				EnvironmentVariables: expandSpringCloudDeploymentEnv(d),
				RuntimeVersion:       appplatform.RuntimeVersion(runtimeVersion),
			},
		},
	}

	future, err := client.Update(ctx, resourceGroup, springCloudName, appName, name, &appResource)
	if err != nil {
		return fmt.Errorf("Error updating Spring Cloud Deployment %q (Spring Cloud Service %q / App Name %q / Resource Group %q): %+v", name, springCloudName, appName, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for update of Spring Cloud Deployment %q (Spring Cloud Service %q / App Name %q / Resource Group %q): %+v", name, springCloudName, appName, resourceGroup, err)
	}

	return resourceArmSpringCloudDeploymentRead(d, meta)
}

func resourceArmSpringCloudDeploymentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).AppPlatform.DeploymentsClient
	appsClient := meta.(*ArmClient).AppPlatform.AppsClient
	ctx, cancel := timeouts.ForDelete(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	springCloudName := id.Path["Spring"]
	appName := id.Path["apps"]
	name := id.Path["deployments"]

	appResp, err := appsClient.Get(ctx, resourceGroup, springCloudName, appName, "")
	if err != nil {
		return fmt.Errorf("Error getting Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q): %+v", appName, springCloudName, resourceGroup, err)
	}

	// active deployment can not be deleted, so just return nil
	if *appResp.Properties.ActiveDeploymentName == name {
		return nil
	}

	if _, err := client.Delete(ctx, resourceGroup, springCloudName, appName, name); err != nil {
		return fmt.Errorf("Error deleting Spring Cloud Deployment %q (Spring Cloud Service %q / App Name %q / Resource Group %q): %+v", name, springCloudName, appName, resourceGroup, err)
	}

	return nil
}

func expandSpringCloudDeploymentUserSourceInfo(d *schema.ResourceData, meta interface{}, resourceGroup, springCloudName, appName string) (*appplatform.UserSourceInfo, error) {
	userSourceInfo := appplatform.UserSourceInfo{
		Type:         appplatform.Jar,
		RelativePath: utils.String("<default>"),
	}
	if jarFile, ok := d.GetOk("jar_file"); ok {
		appsClient := meta.(*ArmClient).AppPlatform.AppsClient
		ctx, cancel := timeouts.ForCreate(meta.(*ArmClient).StopContext, d)
		defer cancel()

		resourceUploadDefinition, err := appsClient.GetResourceUploadURL(ctx, resourceGroup, springCloudName, appName)
		if err != nil {
			return nil, fmt.Errorf("Error getting uploading URL of Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q): %+v", appName, springCloudName, resourceGroup, err)
		}

		if err := uploadJar(d, meta, resourceUploadDefinition, jarFile.(string)); err != nil {
			return nil, fmt.Errorf("Error uploading files %s, %+v", jarFile.(string), err)
		}

		userSourceInfo.RelativePath = resourceUploadDefinition.RelativePath
	}
	return &userSourceInfo, nil
}

func uploadJar(d *schema.ResourceData, meta interface{}, resourceUploadDefinition appplatform.ResourceUploadDefinition, jarFile string) error {
	ctx, cancel := timeouts.ForCreate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	u, err := url.Parse(*resourceUploadDefinition.UploadURL)
	if err != nil {
		return fmt.Errorf("Error parsing uploading URL %s : %+v", *resourceUploadDefinition.UploadURL, err)
	}
	accountName := strings.Split(u.Host, ".")[0]
	shareName := strings.Split(u.Path, "/")[1]
	sasToken := u.RawQuery

	index := strings.LastIndex(*resourceUploadDefinition.RelativePath, "/")
	path := (*resourceUploadDefinition.RelativePath)[:index]
	filename := (*resourceUploadDefinition.RelativePath)[index+1:]

	storageClient := meta.(*ArmClient).Storage
	filesClient := storageClient.FileFilesClientWithSASToken(sasToken)

	file, err := os.Open(jarFile)
	if err != nil {
		return fmt.Errorf("Error opening: %+v", err)
	}

	log.Printf("[DEBUG] Creating Top Level File... accountName: %s shareName: %s, path: %s, filename: %s", accountName, shareName, path, filename)
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("Error 'stat'-ing: %+v", err)
	}
	createFileInput := files.CreateInput{
		ContentLength: info.Size(),
	}
	if _, err := filesClient.Create(ctx, accountName, shareName, path, filename, createFileInput); err != nil {
		return fmt.Errorf("Error creating Top-Level File: %+v", err)
	}
	if err := filesClient.PutFile(ctx, accountName, shareName, path, filename, file, 4); err != nil {
		return fmt.Errorf("Error uploading file : %+v", err)
	}
	return nil
}

func expandSpringCloudDeploymentEnv(d *schema.ResourceData) map[string]*string {
	input := d.Get("env").(map[string]interface{})
	output := make(map[string]*string, len(input))

	for k, v := range input {
		output[k] = utils.String(v.(string))
	}

	return output
}

func flattenSpringCloudDeploymentEnv(input map[string]*string) map[string]string {
	output := make(map[string]string)
	for k, v := range input {
		output[k] = *v
	}

	return output
}
