package managedservices

import (
	"fmt"
	"time"
	"strings"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/managedservices/mgmt/2019-06-01/managedservices"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
				Optional:     true,
				Computed: 	  true,
				ForceNew: 	  true,
				ValidateFunc: validation.IsUUID,
			},

			"scope": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"name": {
				Type:         schema.TypeString,
				Required:     true,
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
			},

			"managed_by_tenant_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},

			"authorization": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"principal_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
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
	client := meta.(*clients.Client).ManagedServices.RegistrationDefinitionsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	registrationDefinitionId := d.Get("registration_definition_id").(string)
	if registrationDefinitionId == "" {
		uuid, err := uuid.GenerateUUID()
		if err != nil {
			return fmt.Errorf("Error generating UUID for Registration Definition: %+v", err)
		}

		registrationDefinitionId = uuid
	}
	
	scope := d.Get("scope").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, scope, registrationDefinitionId)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Registration Definition %q (Scope %q): %+v", registrationDefinitionId, scope, err)
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
			RegistrationDefinitionName: utils.String(d.Get("name").(string)),
			ManagedByTenantID: 			utils.String(d.Get("managed_by_tenant_id").(string)),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, registrationDefinitionId, scope, parameters); err != nil {
		return fmt.Errorf("Error Creating/Updating Registration Definition %q (Scope %q): %+v", registrationDefinitionId, scope, err)
	}

	read, err := client.Get(ctx, scope, registrationDefinitionId)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Registration Definition %q ID (scope %q) ID", registrationDefinitionId, scope)
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
	client := meta.(*clients.Client).ManagedServices.RegistrationDefinitionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parseAzureRegistrationDefinitionId(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.scope, id.registrationDefinitionId)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] Registration Definition '%s' was not found (Scope '%s')", id.registrationDefinitionId, id.scope)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Registration Definition %q (Scope %q): %+v", id.registrationDefinitionId, id.scope, err)
	}

	d.Set("registration_definition_id", resp.Name)
	d.Set("scope", id.scope)

	if props := resp.Properties; props != nil {
		if err := d.Set("authorization", flattenManagedServicesDefinitionAuthorization(props.Authorizations)); err != nil {
			return fmt.Errorf("setting `authorization`: %+v", err)
		}
		d.Set("description", props.Description)
		d.Set("name", props.RegistrationDefinitionName)
		d.Set("managed_by_tenant_id", props.ManagedByTenantID)
	}

	return nil
}

type registrationDefinitionID struct {
	scope 						string
	registrationDefinitionId 	string
}

func parseAzureRegistrationDefinitionId(id string) (*registrationDefinitionID, error) {
	segments := strings.Split(id, "/providers/Microsoft.ManagedServices/registrationDefinitions/")

	if len(segments) != 2 {
		return nil, fmt.Errorf("Expected ID to be in the format `{scope}/providers/Microsoft.ManagedServices/registrationDefinitions/{name} - got %d segments", len(segments))
	}

	azureRegistrationDefinitionId := registrationDefinitionID{
		scope: 						segments[0],
		registrationDefinitionId:  	segments[1],
	}

	return &azureRegistrationDefinitionId, nil
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
	client := meta.(*clients.Client).ManagedServices.RegistrationDefinitionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parseAzureRegistrationDefinitionId(d.Id())
	if err != nil {
		return err
	}
	// The sleep is needed to ensure the registration assignment is successfully deleted 
	// before deleting the registration definition. Bug # is logged with the Product team to track this issue. 
	time.Sleep(20 * time.Second)

	_, err = client.Delete(ctx, id.registrationDefinitionId, id.scope)
	if err != nil {
		return fmt.Errorf("Error deleting Registration Definition %q at Scope %q: %+v", id.registrationDefinitionId, id.scope, err)
	}

	return nil
}
