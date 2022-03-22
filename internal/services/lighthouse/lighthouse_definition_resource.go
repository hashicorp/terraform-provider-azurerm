package lighthouse

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/managedservices/mgmt/2019-06-01/managedservices"
	frsUUID "github.com/gofrs/uuid"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/lighthouse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceLighthouseDefinition() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLighthouseDefinitionCreateUpdate,
		Read:   resourceLighthouseDefinitionRead,
		Update: resourceLighthouseDefinitionCreateUpdate,
		Delete: resourceLighthouseDefinitionDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.LighthouseDefinitionID(id)
			return err
		}),

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
				ValidateFunc: commonids.ValidateSubscriptionID,
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
			return fmt.Errorf("generating UUID for Lighthouse Definition: %+v", err)
		}

		lighthouseDefinitionID = uuid
	}

	id := parse.NewLighthouseDefinitionID(d.Get("scope").(string), lighthouseDefinitionID)
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.Scope, id.LighthouseDefinitionID)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_lighthouse_definition", id.ID())
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

	// NOTE: this API call uses DefinitionId then Scope - check in the future
	if _, err := client.CreateOrUpdate(ctx, id.LighthouseDefinitionID, id.Scope, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
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
			log.Printf("[WARN] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("lighthouse_definition_id", id.LighthouseDefinitionID)
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
		return fmt.Errorf("deleting %s: %+v", *id, err)
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
