package managedapplication

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-06-01/managedapplications"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmManagedApplicationDefinition() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmManagedApplicationDefinitionCreateUpdate,
		Read:   resourceArmManagedApplicationDefinitionRead,
		Update: resourceArmManagedApplicationDefinitionCreateUpdate,
		Delete: resourceArmManagedApplicationDefinitionDelete,

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
				ValidateFunc: ValidateManagedAppDefinitionName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"authorization": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_principal_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.GUID,
						},
						"role_definition_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.GUID,
						},
					},
				},
			},

			"display_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: ValidateManagedAppDefinitionDisplayName,
			},

			"lock_level": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(managedapplications.ReadOnly),
					string(managedapplications.None),
				}, false),
			},

			"create_ui_definition": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.ValidateJsonString,
				DiffSuppressFunc: structure.SuppressJsonDiff,
				ConflictsWith:    []string{"package_file_uri"},
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: ValidateManagedAppDefinitionDescription,
			},

			"main_template": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.ValidateJsonString,
				DiffSuppressFunc: structure.SuppressJsonDiff,
				ConflictsWith:    []string{"package_file_uri"},
			},

			"package_file_uri": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.URLIsHTTPS,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmManagedApplicationDefinitionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagedApplication.ApplicationDefinitionClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroupName := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroupName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for present of existing Managed Application Definition (Managed Application Definition Name %q / Resource Group %q): %+v", name, resourceGroupName, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_managed_application_definition", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	authorizations := d.Get("authorization").(*schema.Set).List()
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	lockLevel := d.Get("lock_level").(string)
	t := d.Get("tags").(map[string]interface{})

	parameters := managedapplications.ApplicationDefinition{
		Location: utils.String(location),
		ApplicationDefinitionProperties: &managedapplications.ApplicationDefinitionProperties{
			Authorizations: expandArmManagedApplicationDefinitionAuthorization(authorizations),
			Description:    utils.String(description),
			DisplayName:    utils.String(displayName),
			LockLevel:      managedapplications.ApplicationLockLevel(lockLevel),
		},
		Tags: tags.Expand(t),
	}

	if v, ok := d.GetOk("create_ui_definition"); ok {
		parameters.CreateUIDefinition = utils.String(v.(string))
	}

	if v, ok := d.GetOk("main_template"); ok {
		parameters.MainTemplate = utils.String(v.(string))
	}

	if v, ok := d.GetOk("package_file_uri"); ok {
		parameters.PackageFileURI = utils.String(v.(string))
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroupName, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating Managed Application Definition (Managed Application Definition Name %q / Resource Group %q): %+v", name, resourceGroupName, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Managed Application Definition (Managed Application Definition Name %q / Resource Group %q): %+v", name, resourceGroupName, err)
	}

	resp, err := client.Get(ctx, resourceGroupName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Managed Application Definition (Managed Application Definition Name %q / Resource Group %q): %+v", name, resourceGroupName, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("Cannot read Managed Application Definition (Managed Application Definition Name %q / Resource Group %q) ID", name, resourceGroupName)
	}
	d.SetId(*resp.ID)

	return resourceArmManagedApplicationDefinitionRead(d, meta)
}

func resourceArmManagedApplicationDefinitionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagedApplication.ApplicationDefinitionClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroupName := id.ResourceGroup
	name := id.Path["applicationDefinitions"]

	resp, err := client.Get(ctx, resourceGroupName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Managed Application Definition %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Managed Application Definition (Managed Application Definition Name %q / Resource Group %q): %+v", name, resourceGroupName, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroupName)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if props := resp.ApplicationDefinitionProperties; props != nil {
		if err := d.Set("authorization", flattenArmManagedApplicationDefinitionAuthorization(props.Authorizations)); err != nil {
			return fmt.Errorf("Error setting `authorization`: %+v", err)
		}
		d.Set("description", props.Description)
		d.Set("display_name", props.DisplayName)
		d.Set("lock_level", string(props.LockLevel))
	}
	if v, ok := d.GetOk("create_ui_definition"); ok {
		d.Set("create_ui_definition", v.(string))
	}
	if v, ok := d.GetOk("main_template"); ok {
		d.Set("main_template", v.(string))
	}
	if v, ok := d.GetOk("package_file_uri"); ok {
		d.Set("package_file_uri", v.(string))
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmManagedApplicationDefinitionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagedApplication.ApplicationDefinitionClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroupName := id.ResourceGroup
	name := id.Path["applicationDefinitions"]

	future, err := client.Delete(ctx, resourceGroupName, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting Managed Application Definition (Managed Application Definition Name %q / Resource Group %q): %+v", name, resourceGroupName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deleting Managed Application Definition (Managed Application Definition Name %q / Resource Group %q): %+v", name, resourceGroupName, err)
		}
	}

	return nil
}

func expandArmManagedApplicationDefinitionAuthorization(input []interface{}) *[]managedapplications.ApplicationProviderAuthorization {
	results := make([]managedapplications.ApplicationProviderAuthorization, 0)
	for _, item := range input {
		v := item.(map[string]interface{})

		result := managedapplications.ApplicationProviderAuthorization{
			RoleDefinitionID: utils.String(v["role_definition_id"].(string)),
			PrincipalID:      utils.String(v["service_principal_id"].(string)),
		}

		results = append(results, result)
	}
	return &results
}

func flattenArmManagedApplicationDefinitionAuthorization(input *[]managedapplications.ApplicationProviderAuthorization) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		servicePrincipalId := ""
		roleDefinitionId := ""

		if item.PrincipalID != nil {
			servicePrincipalId = *item.PrincipalID
		}
		if item.RoleDefinitionID != nil {
			roleDefinitionId = *item.RoleDefinitionID
		}

		results = append(results, map[string]interface{}{
			"role_definition_id":   roleDefinitionId,
			"service_principal_id": servicePrincipalId,
		})
	}

	return results
}
