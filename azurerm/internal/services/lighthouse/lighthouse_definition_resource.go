package lighthouse

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/managedservices/mgmt/2019-06-01/managedservices"
	frsUUID "github.com/gofrs/uuid"
	"github.com/hashicorp/go-uuid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/lighthouse/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/subscription/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceLighthouseDefinition() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLighthouseDefinitionCreateUpdate,
		Read:   resourceLighthouseDefinitionRead,
		Update: resourceLighthouseDefinitionCreateUpdate,
		Delete: resourceLighthouseDefinitionDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"managing_tenant_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"scope": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SubscriptionID,
			},

			"authorization": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"principal_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},

						"role_definition_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},

						"principal_display_name": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"delegated_role_definition_ids": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.IsUUID,
							},
						},
					},
				},
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"lighthouse_definition_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"plan": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"publisher": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"product": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"version": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
		},
	}
}

func resourceLighthouseDefinitionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Lighthouse.DefinitionsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	lighthouseDefinitionID := d.Get("lighthouse_definition_id").(string)
	if lighthouseDefinitionID == "" {
		uuid, err := uuid.GenerateUUID()
		if err != nil {
			return fmt.Errorf("Error generating UUID for Lighthouse Definition: %+v", err)
		}

		lighthouseDefinitionID = uuid
	}

	subscriptionID := meta.(*clients.Client).Account.SubscriptionId
	if subscriptionID == "" {
		return fmt.Errorf("Error reading Subscription for Lighthouse Definition %q", lighthouseDefinitionID)
	}

	scope := d.Get("scope").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, scope, lighthouseDefinitionID)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Lighthouse Definition %q (Scope %q): %+v", lighthouseDefinitionID, scope, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_lighthouse_definition", *existing.ID)
		}
	}
	authorizations, err := expandLighthouseDefinitionAuthorization(d.Get("authorization").(*pluginsdk.Set).List())
	if err != nil {
		return err
	}
	parameters := managedservices.RegistrationDefinition{
		Plan: expandLighthouseDefinitionPlan(d.Get("plan").([]interface{})),
		Properties: &managedservices.RegistrationDefinitionProperties{
			Description:                utils.String(d.Get("description").(string)),
			Authorizations:             authorizations,
			RegistrationDefinitionName: utils.String(d.Get("name").(string)),
			ManagedByTenantID:          utils.String(d.Get("managing_tenant_id").(string)),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, lighthouseDefinitionID, scope, parameters); err != nil {
		return fmt.Errorf("Error Creating/Updating Lighthouse Definition %q (Scope %q): %+v", lighthouseDefinitionID, scope, err)
	}

	read, err := client.Get(ctx, scope, lighthouseDefinitionID)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Lighthouse Definition %q ID (scope %q) ID", lighthouseDefinitionID, scope)
	}

	d.SetId(*read.ID)

	return resourceLighthouseDefinitionRead(d, meta)
}

func resourceLighthouseDefinitionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Lighthouse.DefinitionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LighthouseDefinitionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.Scope, id.LighthouseDefinitionID)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] Lighthouse Definition %q was not found (Scope %q)", id.LighthouseDefinitionID, id.Scope)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Lighthouse Definition %q (Scope %q): %+v", id.LighthouseDefinitionID, id.Scope, err)
	}

	d.Set("lighthouse_definition_id", resp.Name)
	d.Set("scope", id.Scope)

	if err := d.Set("plan", flattenLighthouseDefinitionPlan(resp.Plan)); err != nil {
		return fmt.Errorf("setting `plan`: %+v", err)
	}

	if props := resp.Properties; props != nil {
		if err := d.Set("authorization", flattenLighthouseDefinitionAuthorization(props.Authorizations)); err != nil {
			return fmt.Errorf("setting `authorization`: %+v", err)
		}
		d.Set("description", props.Description)
		d.Set("name", props.RegistrationDefinitionName)
		d.Set("managing_tenant_id", props.ManagedByTenantID)
	}

	return nil
}

func resourceLighthouseDefinitionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Lighthouse.DefinitionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LighthouseDefinitionID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.Delete(ctx, id.LighthouseDefinitionID, id.Scope); err != nil {
		return fmt.Errorf("Error deleting Lighthouse Definition %q at Scope %q: %+v", id.LighthouseDefinitionID, id.Scope, err)
	}

	return nil
}

func flattenLighthouseDefinitionAuthorization(input *[]managedservices.Authorization) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		principalID := ""
		if item.PrincipalID != nil {
			principalID = *item.PrincipalID
		}

		roleDefinitionID := ""
		if item.RoleDefinitionID != nil {
			roleDefinitionID = *item.RoleDefinitionID
		}

		principalIDDisplayName := ""
		if item.PrincipalIDDisplayName != nil {
			principalIDDisplayName = *item.PrincipalIDDisplayName
		}

		results = append(results, map[string]interface{}{
			"role_definition_id":            roleDefinitionID,
			"principal_id":                  principalID,
			"principal_display_name":        principalIDDisplayName,
			"delegated_role_definition_ids": flattenLighthouseDefinitionAuthorizationDelegatedRoleDefinitionIds(item.DelegatedRoleDefinitionIds),
		})
	}

	return results
}

func flattenLighthouseDefinitionAuthorizationDelegatedRoleDefinitionIds(input *[]frsUUID.UUID) []interface{} {
	if input == nil {
		return []interface{}{}
	}
	result := make([]interface{}, 0)
	for _, item := range *input {
		result = append(result, item.String())
	}
	return result
}

func expandLighthouseDefinitionAuthorization(input []interface{}) (*[]managedservices.Authorization, error) {
	results := make([]managedservices.Authorization, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		delegatedRoleDefinitionIds, err := expandLighthouseDefinitionAuthorizationDelegatedRoleDefinitionIds(v["delegated_role_definition_ids"].(*pluginsdk.Set).List())
		if err != nil {
			return nil, err
		}
		result := managedservices.Authorization{
			RoleDefinitionID:           utils.String(v["role_definition_id"].(string)),
			PrincipalID:                utils.String(v["principal_id"].(string)),
			PrincipalIDDisplayName:     utils.String(v["principal_display_name"].(string)),
			DelegatedRoleDefinitionIds: delegatedRoleDefinitionIds,
		}
		results = append(results, result)
	}
	return &results, nil
}

func expandLighthouseDefinitionAuthorizationDelegatedRoleDefinitionIds(input []interface{}) (*[]frsUUID.UUID, error) {
	result := make([]frsUUID.UUID, 0)
	for _, item := range input {
		id, err := frsUUID.FromString(item.(string))
		if err != nil {
			return nil, fmt.Errorf("parsing %q as a UUID: %+v", item, err)
		}
		result = append(result, id)
	}
	return &result, nil
}

func expandLighthouseDefinitionPlan(input []interface{}) *managedservices.Plan {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	raw := input[0].(map[string]interface{})
	return &managedservices.Plan{
		Name:      utils.String(raw["name"].(string)),
		Publisher: utils.String(raw["publisher"].(string)),
		Product:   utils.String(raw["product"].(string)),
		Version:   utils.String(raw["version"].(string)),
	}
}

func flattenLighthouseDefinitionPlan(input *managedservices.Plan) []interface{} {
	if input == nil {
		return []interface{}{}
	}
	var name, publisher, product, version string
	if input.Name != nil {
		name = *input.Name
	}
	if input.Publisher != nil {
		publisher = *input.Publisher
	}
	if input.Product != nil {
		product = *input.Product
	}
	if input.Version != nil {
		version = *input.Version
	}
	return []interface{}{
		map[string]interface{}{
			"name":      name,
			"publisher": publisher,
			"product":   product,
			"version":   version,
		},
	}
}
