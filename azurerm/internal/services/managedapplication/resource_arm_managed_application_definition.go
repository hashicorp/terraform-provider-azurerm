package managedapplication

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-07-01/managedapplications"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/managedapplication/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/managedapplication/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmManagedApplicationDefinition() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmManagedApplicationDefinitionCreateUpdate,
		Read:   resourceArmManagedApplicationDefinitionRead,
		Update: resourceArmManagedApplicationDefinitionCreateUpdate,
		Delete: resourceArmManagedApplicationDefinitionDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ManagedApplicationDefinitionID(id)
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
				ValidateFunc: validate.ManagedApplicationDefinitionName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"lock_level": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(managedapplications.CanNotDelete),
					string(managedapplications.None),
					string(managedapplications.ReadOnly),
				}, false),
			},

			"authorization": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_definition_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},
						"service_principal_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},
					},
				},
			},

			"display_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.ManagedApplicationDefinitionDisplayName,
			},

			"create_ui_definition": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: structure.SuppressJsonDiff,
				ConflictsWith:    []string{"package_file_uri"},
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.ManagedApplicationDefinitionDescription,
			},

			"main_template": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: structure.SuppressJsonDiff,
				ConflictsWith:    []string{"package_file_uri"},
			},

			"package_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"package_file_uri": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmManagedApplicationDefinitionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagedApplication.ApplicationDefinitionClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroupName := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroupName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("failed to check for present of existing Managed Application Definition Name %q (Resource Group %q): %+v", name, resourceGroupName, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_managed_application_definition", *existing.ID)
		}
	}

	parameters := managedapplications.ApplicationDefinition{
		Location: utils.String(azure.NormalizeLocation(d.Get("location"))),
		ApplicationDefinitionProperties: &managedapplications.ApplicationDefinitionProperties{
			Authorizations: expandArmManagedApplicationDefinitionAuthorization(d.Get("authorization").(*schema.Set).List()),
			Description:    utils.String(d.Get("description").(string)),
			DisplayName:    utils.String(d.Get("display_name").(string)),
			IsEnabled:      utils.Bool(d.Get("package_enabled").(bool)),
			LockLevel:      managedapplications.ApplicationLockLevel(d.Get("lock_level").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("create_ui_definition"); ok {
		parameters.CreateUIDefinition = utils.String(v.(string))
	}

	if v, ok := d.GetOk("main_template"); ok {
		parameters.MainTemplate = utils.String(v.(string))
	}

	if (parameters.CreateUIDefinition != nil && parameters.MainTemplate == nil) || (parameters.CreateUIDefinition == nil && parameters.MainTemplate != nil) {
		return fmt.Errorf("`create_ui_definition` and `main_template` should be set or not set together")
	}

	if v, ok := d.GetOk("package_file_uri"); ok {
		parameters.PackageFileURI = utils.String(v.(string))
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroupName, name, parameters)
	if err != nil {
		return fmt.Errorf("failed to create Managed Application Definition %q (Resource Group %q): %+v", name, resourceGroupName, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("failed to wait for creation of Managed Application Definition %q (Resource Group %q): %+v", name, resourceGroupName, err)
	}

	resp, err := client.Get(ctx, resourceGroupName, name)
	if err != nil {
		return fmt.Errorf("failed to retrieve Managed Application Definition %q (Resource Group %q): %+v", name, resourceGroupName, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("cannot read Managed Application Definition %q (Resource Group %q) ID", name, resourceGroupName)
	}
	d.SetId(*resp.ID)

	return resourceArmManagedApplicationDefinitionRead(d, meta)
}

func resourceArmManagedApplicationDefinitionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagedApplication.ApplicationDefinitionClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagedApplicationDefinitionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Managed Application Definition %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to read Managed Application Definition %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if props := resp.ApplicationDefinitionProperties; props != nil {
		if err := d.Set("authorization", flattenArmManagedApplicationDefinitionAuthorization(props.Authorizations)); err != nil {
			return fmt.Errorf("setting `authorization`: %+v", err)
		}
		d.Set("description", props.Description)
		d.Set("display_name", props.DisplayName)
		d.Set("package_enabled", props.IsEnabled)
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

	id, err := parse.ManagedApplicationDefinitionID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("failed to delete Managed Application Definition %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("failed to wait for deleting Managed Application Definition (Managed Application Definition Name %q / Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func expandArmManagedApplicationDefinitionAuthorization(input []interface{}) *[]managedapplications.ApplicationAuthorization {
	results := make([]managedapplications.ApplicationAuthorization, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		result := managedapplications.ApplicationAuthorization{
			RoleDefinitionID: utils.String(v["role_definition_id"].(string)),
			PrincipalID:      utils.String(v["service_principal_id"].(string)),
		}

		results = append(results, result)
	}
	return &results
}

func flattenArmManagedApplicationDefinitionAuthorization(input *[]managedapplications.ApplicationAuthorization) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		servicePrincipalId := ""
		if item.PrincipalID != nil {
			servicePrincipalId = *item.PrincipalID
		}

		roleDefinitionId := ""
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
