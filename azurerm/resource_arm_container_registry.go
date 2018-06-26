package azurerm

import (
	"fmt"
	"log"
	"regexp"

	"strings"

	"github.com/Azure/azure-sdk-for-go/services/containerregistry/mgmt/2017-10-01/containerregistry"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmContainerRegistry() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmContainerRegistryCreate,
		Read:   resourceArmContainerRegistryRead,
		Update: resourceArmContainerRegistryUpdate,
		Delete: resourceArmContainerRegistryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		MigrateState:  resourceAzureRMContainerRegistryMigrateState,
		SchemaVersion: 2,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureRMContainerRegistryName,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"location": locationSchema(),

			"sku": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          string(containerregistry.Classic),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
				ValidateFunc: validation.StringInSlice([]string{
					string(containerregistry.Classic),
					string(containerregistry.Basic),
					string(containerregistry.Standard),
					string(containerregistry.Premium),
				}, true),
			},

			"admin_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"storage_account_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"storage_account": {
				Type:       schema.TypeList,
				Optional:   true,
				Deprecated: "`storage_account` has been replaced by `storage_account_id`.",
				MaxItems:   1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"access_key": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
					},
				},
			},

			"login_server": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"admin_username": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"admin_password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmContainerRegistryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).containerRegistryClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO] preparing arguments for AzureRM Container Registry creation.")

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	sku := d.Get("sku").(string)
	adminUserEnabled := d.Get("admin_enabled").(bool)
	tags := d.Get("tags").(map[string]interface{})

	parameters := containerregistry.Registry{
		Location: &location,
		Sku: &containerregistry.Sku{
			Name: containerregistry.SkuName(sku),
			Tier: containerregistry.SkuTier(sku),
		},
		RegistryProperties: &containerregistry.RegistryProperties{
			AdminUserEnabled: utils.Bool(adminUserEnabled),
		},
		Tags: expandTags(tags),
	}

	if v, ok := d.GetOk("storage_account_id"); ok {
		if strings.ToLower(sku) != strings.ToLower(string(containerregistry.Classic)) {
			return fmt.Errorf("`storage_account_id` can only be specified for a Classic (unmanaged) Sku.")
		}

		parameters.StorageAccount = &containerregistry.StorageAccountProperties{
			ID: utils.String(v.(string)),
		}
	} else {
		if strings.ToLower(sku) == strings.ToLower(string(containerregistry.Classic)) {
			return fmt.Errorf("`storage_account_id` must be specified for a Classic (unmanaged) Sku.")
		}
	}

	future, err := client.Create(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating Container Registry %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for creation of Container Registry %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Container Registry %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Container Registry %q (resource group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmContainerRegistryRead(d, meta)
}

func resourceArmContainerRegistryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).containerRegistryClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO] preparing arguments for AzureRM Container Registry update.")

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	sku := d.Get("sku").(string)
	adminUserEnabled := d.Get("admin_enabled").(bool)
	tags := d.Get("tags").(map[string]interface{})

	parameters := containerregistry.RegistryUpdateParameters{
		RegistryPropertiesUpdateParameters: &containerregistry.RegistryPropertiesUpdateParameters{
			AdminUserEnabled: utils.Bool(adminUserEnabled),
		},
		Sku: &containerregistry.Sku{
			Name: containerregistry.SkuName(sku),
			Tier: containerregistry.SkuTier(sku),
		},
		Tags: expandTags(tags),
	}

	if v, ok := d.GetOk("storage_account_id"); ok {
		if strings.ToLower(sku) != strings.ToLower(string(containerregistry.Classic)) {
			return fmt.Errorf("`storage_account_id` can only be specified for a Classic (unmanaged) Sku.")
		}

		parameters.StorageAccount = &containerregistry.StorageAccountProperties{
			ID: utils.String(v.(string)),
		}
	} else {
		if strings.ToLower(sku) == strings.ToLower(string(containerregistry.Classic)) {
			return fmt.Errorf("`storage_account_id` must be specified for a Classic (unmanaged) Sku.")
		}
	}

	future, err := client.Update(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error updating Container Registry %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for update of Container Registry %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Container Registry %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Container Registry %q (resource group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmContainerRegistryRead(d, meta)
}

func resourceArmContainerRegistryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).containerRegistryClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["registries"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Container Registry %q was not found in Resource Group %q", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure Container Registry %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}
	d.Set("admin_enabled", resp.AdminUserEnabled)
	d.Set("login_server", resp.LoginServer)

	if sku := resp.Sku; sku != nil {
		d.Set("sku", string(sku.Tier))
	}

	if account := resp.StorageAccount; account != nil {
		d.Set("storage_account_id", account.ID)
	}

	if *resp.AdminUserEnabled {
		credsResp, err := client.ListCredentials(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Error making Read request on Azure Container Registry %s for Credentials: %s", name, err)
		}

		d.Set("admin_username", credsResp.Username)
		for _, v := range *credsResp.Passwords {
			d.Set("admin_password", v.Value)
			break
		}
	} else {
		d.Set("admin_username", "")
		d.Set("admin_password", "")
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmContainerRegistryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).containerRegistryClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["registries"]

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error issuing Azure ARM delete request of Container Registry '%s': %+v", name, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error issuing Azure ARM delete request of Container Registry '%s': %+v", name, err)
	}

	return nil
}

func validateAzureRMContainerRegistryName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"alpha numeric characters only are allowed in %q: %q", k, value))
	}

	if 5 > len(value) {
		errors = append(errors, fmt.Errorf("%q cannot be less than 5 characters: %q", k, value))
	}

	if len(value) >= 50 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 50 characters: %q %d", k, value, len(value)))
	}

	return
}
