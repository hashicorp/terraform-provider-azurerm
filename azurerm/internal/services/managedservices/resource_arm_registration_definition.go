package managedservices

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/managedservices/mgmt/2019-06-01/managedservices"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmRegistrationDefinition() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmRegistrationDefinitionCreateUpdate,
		Read:   resourceArmRegistrationDefinitionRead,
		Update: resourceArmRegistrationDefinitionCreateUpdate,
		Delete: resourceArmRegistrationDefinitionDelete,
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
			"registration_definition_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(3, 128),
			},

			"registration_definition_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(3, 128),
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(3, 128),
			},

			"scope": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(3, 128),
			},

			"managed_by_tenant_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},

			"managed_by_tenant_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(3, 128),
			},

			"authorizations": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"principal_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},
						"principal_display_name": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(3, 128),
						},
						"role_definition_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},
					},
				},
			},
		},
	}
}

func resourceArmRegistrationDefinitionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagedService.RegistrationDefinitionsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	registrationDefinitionID := d.Get("registration_definition_id").(string)
	scope := d.Get("scope").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, scope, registrationDefinitionID)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Registration Definition %q (Scope %q): %+v", registrationDefinitionID, scope, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_registration_definition", *existing.ID)
		}
	}

	parameters := managedservices.RegistrationDefinition {
		Properties: &managedservices.RegistrationDefinitionProperties{
			Description: 				utils.String(d.Get("description").(string)),
			Authorizations: 			expandManagedServicesDefinitionAuthorization(d.Get("authorization").(*schema.Set).List()),
			RegistrationDefinitionName: utils.String(d.Get("registration_definition_name").(string)),
			ManagedByTenantID: 			utils.String(d.Get("managed_by_tenant_id").(string)),
			ManagedByTenantName:		utils.String(d.Get("managed_by_tenant_name").(string)),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, registrationDefinitionID, scope, parameters); err != nil {
		return fmt.Errorf("Error Creating/Updating Registration Definition %q (Scope %q): %+v", registrationDefinitionID, scope, err)
	}

	read, err := client.Get(ctx, scope, registrationDefinitionID)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Registration Definition %q ID (scope %q) ID", registrationDefinitionID, scope)
	}

	d.SetId(*read.ID)

	return resourceArmRegistrationDefinitionRead(d, meta)
}

func expandManagedServicesDefinitionAuthorization(input []interface{}) *[]managedservices.Authorization {
	results := make([]managedservices.Authorization, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		result := managedservices.Authorization{
			RoleDefinitionID: utils.String(v["role_definition_id"].(string)),
			PrincipalID:      utils.String(v["principal_id"].(string)),
		}
		results = append(results, result)
	}
	return &results
}

func resourceArmRegistrationDefinitionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagedService.RegistrationDefinitionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	registrationDefinitionID:= id.Path["registrationDefinitionID"]
	scope:= id.Path["scope"]

	resp, err := client.Get(ctx, scope, registrationDefinitionID)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Registration Assignment %q (Scope %q): %+v", registrationDefinitionID, scope, err)
	}

	d.Set("registration_definition_id", resp.Name)
	d.Set("scope", scope)

	if props := resp.Properties; props != nil {
		if err := d.Set("authorization", flattenManagedServicesDefinitionAuthorization(props.Authorizations)); err != nil {
			return fmt.Errorf("setting `authorization`: %+v", err)
		}
		d.Set("description", props.Description)
		d.Set("registration_definition_name", props.RegistrationDefinitionName)
		d.Set("managed_by_tenant_id", props.ManagedByTenantID)
		d.Set("managed_by_tenant_name", props.ManagedByTenantName)
	}

	return nil
}

func flattenManagedServicesDefinitionAuthorization(input *[]managedservices.Authorization) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		principalId := ""
		if item.PrincipalID != nil {
			principalId = *item.PrincipalID
		}

		roleDefinitionId := ""
		if item.RoleDefinitionID != nil {
			roleDefinitionId = *item.RoleDefinitionID
		}

		results = append(results, map[string]interface{}{
			"role_definition_id":   roleDefinitionId,
			"principal_id": principalId,
		})
	}

	return results
}

func resourceArmRegistrationDefinitionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagedService.RegistrationDefinitionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	registrationDefinitionID:= id.Path["registrationDefinitionID"]
	scope:= id.Path["scope"]

	_, err = client.Delete(ctx, registrationDefinitionID, scope)
	if err != nil {
		return fmt.Errorf("Error deleting Registration Definition %q (Scope %q): %+v", registrationDefinitionID, scope, err)
	}

	return nil
}
