package blueprint

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2018-09-01-preview/authorization"
	"github.com/Azure/azure-sdk-for-go/services/preview/blueprint/mgmt/2018-11-01-preview/blueprint"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/set"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/blueprint/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/blueprint/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmBlueprintAssignment() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmBlueprintAssignmentCreateUpdate,
		Read:   resourceArmBlueprintAssignmentRead,
		Update: resourceArmBlueprintAssignmentCreateUpdate,
		Delete: resourceArmBlueprintAssignmentDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.BlueprintAssignmentID(id)
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
				ValidateFunc: validate.BlueprintAssignmentName,
			},

			"location": azure.SchemaLocation(),

			// The scope of the resource. Valid scopes are:
			// management group (format: '/providers/Microsoft.Management/managementGroups/{managementGroup}')
			// subscription (format: '/subscriptions/{subscriptionId}')
			"scope": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.BlueprintAssignmentScopeID,
			},

			"blueprint_definition_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.BlueprintDefinitionID,
			},

			"identity": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								// None is listed as possible values in go SDK, but the service will reject type none of identity
								string(blueprint.ManagedServiceIdentityTypeSystemAssigned),
								string(blueprint.ManagedServiceIdentityTypeUserAssigned),
							}, false),
							// The first character of value returned by the service is always in lower case.
							DiffSuppressFunc: suppress.CaseDifference,
						},

						"identity_ids": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								// TODO: validation for a UAI which requires an ID Parser/Validator
							},
						},

						"principal_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"parameter_values": {
				Type:     schema.TypeString,
				Optional: true,
				// This state function is used to normalize the format of the input JSON string,
				// and strip any extra field comparing to the allowed fields in the swagger
				// to avoid unnecessary diff in the state and config
				StateFunc:    normalizeAssignmentParameterValuesJSON,
				ValidateFunc: validation.StringIsJSON,
				// Suppress the differences caused by JSON formatting or ordering
				DiffSuppressFunc: structure.SuppressJsonDiff,
			},

			"resource_groups": {
				Type:     schema.TypeString,
				Optional: true,
				// This state function is used to normalize the format of the input JSON string,
				// and strip any extra field comparing to the allowed fields in the swagger
				// to avoid unnecessary diff in the state and config
				StateFunc:    normalizeAssignmentResourceGroupValuesJSON,
				ValidateFunc: validation.StringIsJSON,
				// Suppress the differences caused by JSON formatting or ordering
				DiffSuppressFunc: structure.SuppressJsonDiff,
			},

			"lock_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(blueprint.None),
				ValidateFunc: validation.StringInSlice([]string{
					string(blueprint.None),
					string(blueprint.AllResourcesReadOnly),
					string(blueprint.AllResourcesDoNotDelete),
				}, false),
				// The first character of value returned by the service is always in lower case.
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"lock_exclude_principals": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 5,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},

			"published_blueprint_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmBlueprintAssignmentCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Blueprint.AssignmentClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	scope := d.Get("scope").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, scope, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("unable to check for presence of existing Blueprint Assignment %q (Scope %q): %+v", name, scope, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_blueprint_assignment", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))

	blueprintID := d.Get("blueprint_definition_id").(string)

	identityRaw := d.Get("identity").([]interface{})
	identity, err := expandArmBlueprintAssignmentIdentity(identityRaw)
	if err != nil {
		return fmt.Errorf("unable to expand `identity`: %+v", err)
	}

	lockMode := d.Get("lock_mode").(string)

	excludedPrincipalsRaw := d.Get("lock_exclude_principals").(*schema.Set)
	excludedPrincipals := utils.ExpandStringSlice(excludedPrincipalsRaw.List())

	assignment := blueprint.Assignment{
		Location: utils.String(location),
		Identity: identity,
		AssignmentProperties: &blueprint.AssignmentProperties{
			BlueprintID: &blueprintID,
			Locks: &blueprint.AssignmentLockSettings{
				Mode:               blueprint.AssignmentLockMode(lockMode),
				ExcludedPrincipals: excludedPrincipals,
			},
		},
	}

	if v, ok := d.GetOk("description"); ok {
		assignment.Description = utils.String(v.(string))
	}

	if v, ok := d.GetOk("display_name"); ok {
		assignment.DisplayName = utils.String(v.(string))
	}

	if v, ok := d.GetOk("parameter_values"); ok {
		assignment.Parameters = expandArmBlueprintAssignmentParameters(v.(string))
	}

	if v, ok := d.GetOk("resource_groups"); ok {
		assignment.ResourceGroups = expandArmBlueprintAssignmentResourceGroups(v.(string))
	}

	// if the identity of blueprint assignment is SystemAssigned, owner permission needs to be granted
	if assignment.Identity.Type == blueprint.ManagedServiceIdentityTypeSystemAssigned {
		log.Printf("[DEBUG] Need to grant owner permission for blueprint assignment when identity is set to SystemAssigned")
		// get SPN object ID of the blueprint
		resp, err := client.WhoIsBlueprint(ctx, scope, name)
		if err != nil {
			return fmt.Errorf("unable to create Blueprint Assignment %q (Scope %q) with SystemAssigned identity: %+v", name, scope, err)
		}
		if resp.ObjectID == nil {
			return fmt.Errorf("unable to create Blueprint Assignment %q (Scope %q) with SystemAssigned identity: The SPN Object ID of the assignment is nil", name, scope)
		}

		// get owner role definition id
		log.Printf("[DEBUG] query the ID of the owner role definition")
		roleDefinitionClient := meta.(*clients.Client).Authorization.RoleDefinitionsClient
		roleDefinitions, err := roleDefinitionClient.List(ctx, scope, "roleName eq 'owner'")
		if err != nil {
			return fmt.Errorf("unable to create Blueprint Assignment %q (Scope %q) with SystemAssigned identity: %+v", name, scope, err)
		}
		if len(roleDefinitions.Values()) != 1 {
			return fmt.Errorf("unable to create Blueprint Assignment %q (Scope %q) with SystemAssigned identity: cannot load the owner role definition id", name, scope)
		}
		ownerDefinitionId := roleDefinitions.Values()[0].ID
		log.Printf("[DEBUG] owner definition ID: %s", *ownerDefinitionId)

		// assign owner permission to blueprint assignment
		roleAssignmentClient := meta.(*clients.Client).Authorization.RoleAssignmentsClient
		// query if ownership assignment already exists
		log.Printf("[DEBUG] querying whether is owner permission has already assigned")
		result, err := roleAssignmentClient.List(ctx, fmt.Sprintf("principalId eq '%s'", *resp.ObjectID))
		if err != nil {
			return fmt.Errorf("unable to create Blueprint Assignment %q (Scope %q) with SystemAssigned identity: %+v", name, scope, err)
		}
		ownerPermissionAssigned := false
		for _, v := range result.Values() {
			if v.RoleAssignmentPropertiesWithScope == nil || v.RoleAssignmentPropertiesWithScope.RoleDefinitionID == nil {
				continue
			}
			if *v.RoleAssignmentPropertiesWithScope.RoleDefinitionID == *ownerDefinitionId {
				ownerPermissionAssigned = true
				break
			}
		}
		log.Printf("[DEBUG] owner permission assigned: %v", ownerPermissionAssigned)
		if !ownerPermissionAssigned {
			// assign owner permission
			roleAssignmentName, err := uuid.GenerateUUID()
			if err != nil {
				return fmt.Errorf("unable to create Blueprint Assignment %q (Scope %q) with SystemAssigned identity: %+v", name, scope, err)
			}
			log.Printf("[DEBUG] assigning owner role for blueprint with uuid '%s'", roleAssignmentName)
			roleAssignmentParameters := authorization.RoleAssignmentCreateParameters{
				RoleAssignmentProperties: &authorization.RoleAssignmentProperties{
					RoleDefinitionID: ownerDefinitionId,
					PrincipalID:      resp.ObjectID,
				},
			}
			_, err = roleAssignmentClient.Create(ctx, scope, roleAssignmentName, roleAssignmentParameters)
			if err != nil {
				return fmt.Errorf("unable to create Blueprint Assignment %q (Scope %q) with SystemAssigned identity: %+v", roleAssignmentName, scope, err)
			}
		}
	}

	if _, err := client.CreateOrUpdate(ctx, scope, name, assignment); err != nil {
		return fmt.Errorf("unable to create Blueprint Assignment %q (Scope %q): %+v", name, scope, err)
	}

	// the blueprint assignment is not ready after creation until its provisioning state turns to "succeeded"
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			string(blueprint.Waiting),
			string(blueprint.Validating),
			string(blueprint.Creating),
			string(blueprint.Deploying),
			string(blueprint.Locking),
		},
		Target:  []string{string(blueprint.Succeeded)},
		Refresh: blueprintAssignmentCreateStateRefreshFunc(ctx, client, scope, name),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Failed waiting for Blueprint Assignment %q (Scope %q): %+v", name, scope, err)
	}

	resp, err := client.Get(ctx, scope, name)
	if err != nil {
		return fmt.Errorf("unable to retrieve Blueprint Assignment %q (Scope %q): %+v", name, scope, err)
	}

	d.SetId(*resp.ID)

	return resourceArmBlueprintAssignmentRead(d, meta)
}

func resourceArmBlueprintAssignmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Blueprint.AssignmentClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BlueprintAssignmentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ScopeId, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Blueprint Assignment %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("unable to read Blueprint Assignment %q (Scope %q): %+v", id.Name, id.ScopeId, err)
	}

	d.Set("name", id.Name)
	d.Set("scope", id.ScopeId)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if err := d.Set("identity", flattenArmBlueprintAssignmentIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("unable to set `identity`: %+v", err)
	}

	if resp.AssignmentProperties == nil {
		return fmt.Errorf("unable to retrieve Blueprint Assignment %q (Scope %q): `properties` was nil", id.Name, id.ScopeId)
	}

	props := *resp.AssignmentProperties
	// the `BlueprintID` field in response is the ID of the published blueprint (with version) which is different with the blueprint ID in user's input
	if props.BlueprintID == nil {
		return fmt.Errorf("unable to retrieve Blueprint Assignment %q (Scope %q): BlueprintID is nil", id.Name, id.ScopeId)
	}

	publishedID, err := parse.PublishedBlueprintID(*props.BlueprintID)
	if err != nil {
		return fmt.Errorf("unable to retrieve Blueprint Assignment %q (Scope %q): %+v", id.Name, id.ScopeId, err)
	}

	d.Set("blueprint_definition_id", publishedID.BlueprintDefinitionId.ID)
	d.Set("published_blueprint_id", props.BlueprintID)
	d.Set("description", props.Description)
	d.Set("display_name", props.DisplayName)

	// flatten parameter_values
	paramValues, err := flattenArmBlueprintAssignmentParameters(props.Parameters)
	if err != nil {
		return fmt.Errorf("unable to flatten `parameter_values`: %+v", err)
	}
	if err := d.Set("parameter_values", paramValues); err != nil {
		return fmt.Errorf("unable to set `parameter_values`: %+v", err)
	}

	// flatten resource_groups
	resGroups, err := flattenArmBlueprintAssignmentResourceGroups(props.ResourceGroups)
	if err != nil {
		return fmt.Errorf("unable to flatten `resource_groups`: %+v", err)
	}
	if err := d.Set("resource_groups", resGroups); err != nil {
		return fmt.Errorf("unable to set `resource_groups`: %+v", err)
	}

	lockMode := string(blueprint.None)
	excludePrincipals := new(schema.Set)
	if lock := props.Locks; lock != nil {
		lockMode = string(lock.Mode)
		if lock.ExcludedPrincipals != nil {
			excludePrincipals = set.FromStringSlice(*lock.ExcludedPrincipals)
		}
	}
	d.Set("lock_mode", lockMode)
	d.Set("lock_exclude_principals", excludePrincipals)

	return nil
}

func resourceArmBlueprintAssignmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Blueprint.AssignmentClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BlueprintAssignmentID(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Delete(ctx, id.ScopeId, id.Name)
	if err != nil {
		return fmt.Errorf("unable to delete Blueprint Assignment %q (Scope %q): %+v", id.Name, id.ScopeId, err)
	}

	// the blueprint assignment is not deleted immediately after the Delete func returns.
	// There is a short period that the provisioning state will turn to "deleting" before it is really deleted.
	stateConf := &resource.StateChangeConf{
		Pending: []string{"200"},
		Target:  []string{"404"},
		Refresh: blueprintAssignmentDeleteStateRefreshFunc(ctx, client, id.ScopeId, id.Name),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Failed waiting for the deletion of Blueprint Assignment %q (Scope %q): %+v", id.Name, id.ScopeId, err)
	}

	return nil
}

func blueprintAssignmentCreateStateRefreshFunc(ctx context.Context, client *blueprint.AssignmentsClient, scope, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, scope, name)
		if err != nil {
			return nil, "", fmt.Errorf("unable to retrieve Blueprint Assignment %q (Scope %q): %+v", name, scope, err)
		}

		return resp, string(resp.ProvisioningState), nil
	}
}

func blueprintAssignmentDeleteStateRefreshFunc(ctx context.Context, client *blueprint.AssignmentsClient, scope, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, scope, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return nil, "", fmt.Errorf("unable to retrieve Blueprint Assignment %q (Scope %q): %+v", name, scope, err)
			}
		}

		return resp, strconv.Itoa(resp.StatusCode), nil
	}
}

func expandArmBlueprintAssignmentIdentity(input []interface{}) (*blueprint.ManagedServiceIdentity, error) {
	if len(input) == 0 {
		return nil, nil
	}

	raw := input[0].(map[string]interface{})

	identity := blueprint.ManagedServiceIdentity{
		Type: blueprint.ManagedServiceIdentityType(raw["type"].(string)),
	}

	identityIdsRaw := raw["identity_ids"].(*schema.Set).List()
	identityIds := make(map[string]*blueprint.UserAssignedIdentity)
	for _, v := range identityIdsRaw {
		identityIds[v.(string)] = &blueprint.UserAssignedIdentity{}
	}

	if len(identityIds) > 0 {
		if identity.Type != blueprint.ManagedServiceIdentityTypeUserAssigned {
			return nil, fmt.Errorf("`identity_ids` can only be specified when `type` includes `UserAssigned`")
		}

		identity.UserAssignedIdentities = identityIds
	}

	return &identity, nil
}

func flattenArmBlueprintAssignmentIdentity(input *blueprint.ManagedServiceIdentity) []interface{} {
	if input == nil || input.Type == blueprint.ManagedServiceIdentityTypeNone {
		return []interface{}{}
	}

	identityIds := make([]string, 0)
	if input.UserAssignedIdentities != nil {
		for k := range input.UserAssignedIdentities {
			identityIds = append(identityIds, k)
		}
	}

	principalId := ""
	if input.PrincipalID != nil {
		principalId = *input.PrincipalID
	}

	tenantId := ""
	if input.TenantID != nil {
		tenantId = *input.TenantID
	}

	return []interface{}{
		map[string]interface{}{
			"type":         string(input.Type),
			"identity_ids": identityIds,
			"principal_id": principalId,
			"tenant_id":    tenantId,
		},
	}
}

func expandArmBlueprintAssignmentParameters(input string) map[string]*blueprint.ParameterValue {
	var result map[string]*blueprint.ParameterValue
	// the string has been validated by the schema, therefore the error is ignored here, since it will never happen.
	_ = json.Unmarshal([]byte(input), &result)
	return result
}

func flattenArmBlueprintAssignmentParameters(input map[string]*blueprint.ParameterValue) (string, error) {
	if len(input) == 0 {
		return "", nil
	}

	b, err := json.Marshal(input)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func expandArmBlueprintAssignmentResourceGroups(input string) map[string]*blueprint.ResourceGroupValue {
	var result map[string]*blueprint.ResourceGroupValue
	// the string has been validated by the schema, therefore the error is ignored here, since it will never happen.
	_ = json.Unmarshal([]byte(input), &result)
	return result
}

func flattenArmBlueprintAssignmentResourceGroups(input map[string]*blueprint.ResourceGroupValue) (string, error) {
	if len(input) == 0 {
		return "", nil
	}

	b, err := json.Marshal(input)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func normalizeAssignmentParameterValuesJSON(jsonString interface{}) string {
	if jsonString == nil || jsonString == "" {
		return ""
	}

	var values map[string]*blueprint.ParameterValue
	if err := json.Unmarshal([]byte(jsonString.(string)), &values); err != nil {
		return fmt.Sprintf("unable to parse JSON: %+v", err)
	}

	b, _ := json.Marshal(values)
	return string(b)
}

func normalizeAssignmentResourceGroupValuesJSON(jsonString interface{}) string {
	if jsonString == nil || jsonString == "" {
		return ""
	}

	var values map[string]*blueprint.ResourceGroupValue
	if err := json.Unmarshal([]byte(jsonString.(string)), &values); err != nil {
		return fmt.Sprintf("unable to parse JSON: %+v", err)
	}

	b, _ := json.Marshal(values)
	return string(b)
}
