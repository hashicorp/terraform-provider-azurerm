package authorization

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2020-04-01-preview/authorization"
	"github.com/hashicorp/go-uuid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/authorization/azuresdkhacks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/authorization/migration"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/authorization/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmRoleDefinition() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmRoleDefinitionCreate,
		Read:   resourceArmRoleDefinitionRead,
		Update: resourceArmRoleDefinitionUpdate,
		Delete: resourceArmRoleDefinitionDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.RoleDefinitionId(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.RoleDefinitionV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"role_definition_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"scope": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			//lintignore:XS003
			"permissions": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"actions": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
						"not_actions": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
						"data_actions": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
							Set: pluginsdk.HashString,
						},
						"not_data_actions": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
							Set: pluginsdk.HashString,
						},
					},
				},
			},

			"assignable_scopes": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"role_definition_resource_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmRoleDefinitionCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Authorization.RoleDefinitionsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	roleDefinitionId := d.Get("role_definition_id").(string)
	if roleDefinitionId == "" {
		uuid, err := uuid.GenerateUUID()
		if err != nil {
			return fmt.Errorf("generating UUID for Role Assignment: %+v", err)
		}

		roleDefinitionId = uuid
	}

	name := d.Get("name").(string)
	scope := d.Get("scope").(string)
	description := d.Get("description").(string)
	roleType := "CustomRole"

	permissionsRaw := d.Get("permissions").([]interface{})
	permissions := expandRoleDefinitionPermissions(permissionsRaw)
	assignableScopes := expandRoleDefinitionAssignableScopes(d)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, scope, roleDefinitionId)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Role Definition ID for %q (Scope %q)", name, scope)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			importID := fmt.Sprintf("%s|%s", *existing.ID, scope)
			return tf.ImportAsExistsError("azurerm_role_definition", importID)
		}
	}

	properties := authorization.RoleDefinition{
		RoleDefinitionProperties: &authorization.RoleDefinitionProperties{
			RoleName:         utils.String(name),
			Description:      utils.String(description),
			RoleType:         utils.String(roleType),
			Permissions:      &permissions,
			AssignableScopes: &assignableScopes,
		},
	}

	if _, err := client.CreateOrUpdate(ctx, scope, roleDefinitionId, properties); err != nil {
		return err
	}

	// (@jackofallops) - Updates are subject to eventual consistency, and could be read as stale data
	if !d.IsNewResource() {
		id, err := parse.RoleDefinitionId(d.Id())
		if err != nil {
			return err
		}
		stateConf := &pluginsdk.StateChangeConf{
			Pending: []string{
				"Pending",
			},
			Target: []string{
				"OK",
			},
			Refresh:                   roleDefinitionUpdateStateRefreshFunc(ctx, client, id.ResourceID),
			MinTimeout:                10 * time.Second,
			ContinuousTargetOccurence: 12,
			Timeout:                   d.Timeout(pluginsdk.TimeoutUpdate),
		}

		if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for update to Role Definition %q to finish replicating", name)
		}
	}

	read, err := client.Get(ctx, scope, roleDefinitionId)
	if err != nil {
		return err
	}
	if read.ID == nil || *read.ID == "" {
		return fmt.Errorf("Cannot read Role Definition ID for %q (Scope %q)", name, scope)
	}

	d.SetId(fmt.Sprintf("%s|%s", *read.ID, scope))
	return resourceArmRoleDefinitionRead(d, meta)
}

func resourceArmRoleDefinitionUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	sdkClient := meta.(*clients.Client).Authorization.RoleDefinitionsClient
	client := azuresdkhacks.NewRoleDefinitionsWorkaroundClient(sdkClient)
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	roleDefinitionId, err := parse.RoleDefinitionId(d.Id())
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	roleType := "CustomRole"

	permissionsRaw := d.Get("permissions").([]interface{})
	permissions := expandRoleDefinitionPermissions(permissionsRaw)
	assignableScopes := expandRoleDefinitionAssignableScopes(d)

	properties := authorization.RoleDefinition{
		RoleDefinitionProperties: &authorization.RoleDefinitionProperties{
			RoleName:         utils.String(name),
			Description:      utils.String(description),
			RoleType:         utils.String(roleType),
			Permissions:      &permissions,
			AssignableScopes: &assignableScopes,
		},
	}

	resp, err := client.CreateOrUpdate(ctx, roleDefinitionId.Scope, roleDefinitionId.RoleID, properties)
	if err != nil {
		return fmt.Errorf("updating Role Definition %q (Scope %q): %+v", roleDefinitionId.RoleID, roleDefinitionId.Scope, err)
	}
	if resp.RoleDefinitionProperties == nil {
		return fmt.Errorf("updating Role Definition %q (Scope %q): `properties` was nil", roleDefinitionId.RoleID, roleDefinitionId.Scope)
	}
	updatedOn := resp.RoleDefinitionProperties.UpdatedOn
	if updatedOn == nil {
		return fmt.Errorf("updating Role Definition %q (Scope %q): `properties.UpdatedOn` was nil", roleDefinitionId.RoleID, roleDefinitionId.Scope)
	}

	// "Updating" a role definition actually creates a new one and these get consolidated a few seconds later
	// where the "create date" and "update date" match for the newly created record
	// but eventually switch to being the old create date and the new update date
	// ergo we can can for the old create date and the new updated date
	log.Printf("[DEBUG] Waiting for Role Definition %q (Scope %q) to settle down..", roleDefinitionId.RoleID, roleDefinitionId.Scope)
	stateConf := &pluginsdk.StateChangeConf{
		ContinuousTargetOccurence: 12,
		Delay:                     60 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{"Pending"},
		Target:                    []string{"Updated"},
		Refresh:                   roleDefinitionEventualConsistencyUpdate(ctx, client, *roleDefinitionId, *updatedOn),
		Timeout:                   d.Timeout(pluginsdk.TimeoutUpdate),
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for Role Definition %q (Scope %q) to settle down: %+v", roleDefinitionId.RoleID, roleDefinitionId.Scope, err)
	}

	return resourceArmRoleDefinitionRead(d, meta)
}

func resourceArmRoleDefinitionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Authorization.RoleDefinitionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	roleDefinitionId, err := parse.RoleDefinitionId(d.Id())
	if err != nil {
		return err
	}

	d.Set("scope", roleDefinitionId.Scope)
	d.Set("role_definition_id", roleDefinitionId.RoleID)
	d.Set("role_definition_resource_id", roleDefinitionId.ResourceID)

	resp, err := client.Get(ctx, roleDefinitionId.Scope, roleDefinitionId.RoleID)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Role Definition %q was not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("loading Role Definition %q: %+v", d.Id(), err)
	}

	if props := resp.RoleDefinitionProperties; props != nil {
		d.Set("name", props.RoleName)
		d.Set("description", props.Description)

		permissions := flattenRoleDefinitionPermissions(props.Permissions)
		if err := d.Set("permissions", permissions); err != nil {
			return err
		}

		assignableScopes := flattenRoleDefinitionAssignableScopes(props.AssignableScopes)
		if err := d.Set("assignable_scopes", assignableScopes); err != nil {
			return err
		}
	}

	return nil
}

func resourceArmRoleDefinitionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Authorization.RoleDefinitionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, _ := parse.RoleDefinitionId(d.Id())

	resp, err := client.Delete(ctx, id.Scope, id.RoleID)
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("deleting Role Definition %q at Scope %q: %+v", id.RoleID, id.Scope, err)
		}
	}
	// Deletes are not instant and can take time to propagate
	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{
			"Pending",
		},
		Target: []string{
			"Deleted",
			"NotFound",
		},
		Refresh:                   roleDefinitionDeleteStateRefreshFunc(ctx, client, id.ResourceID),
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 6,
		Timeout:                   d.Timeout(pluginsdk.TimeoutDelete),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for delete on Role Definition %q to complete", id.RoleID)
	}

	return nil
}

func roleDefinitionEventualConsistencyUpdate(ctx context.Context, client azuresdkhacks.RoleDefinitionsWorkaroundClient, id parse.RoleDefinitionID, updateRequestDate string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, id.Scope, id.RoleID)
		if err != nil {
			return resp, "Failed", err
		}
		if resp.RoleDefinitionProperties == nil {
			return resp, "Failed", fmt.Errorf("`properties` was nil")
		}
		if resp.RoleDefinitionProperties.CreatedOn == nil {
			return resp, "Failed", fmt.Errorf("`properties.CreatedOn` was nil")
		}

		if resp.RoleDefinitionProperties.UpdatedOn == nil {
			return resp, "Failed", fmt.Errorf("`properties.UpdatedOn` was nil")
		}

		updateRequestTime, err := time.Parse(time.RFC3339, updateRequestDate)
		if err != nil {
			return nil, "", fmt.Errorf("error parsing time from update request: %+v", err)
		}

		respCreatedOn, err := time.Parse(time.RFC3339, *resp.RoleDefinitionProperties.CreatedOn)
		if err != nil {
			return nil, "", fmt.Errorf("error parsing time for createdOn from update request: %+v", err)
		}

		respUpdatedOn, err := time.Parse(time.RFC3339, *resp.RoleDefinitionProperties.UpdatedOn)
		if err != nil {
			return nil, "", fmt.Errorf("error parsing time for updatedOn from update request: %+v", err)
		}

		if respCreatedOn.Equal(updateRequestTime) {
			// a new role definition is created and eventually (~5s) reconciled
			return resp, "Pending", nil
		}

		if !respUpdatedOn.After(updateRequestTime) {
			// The real updated on will be after the time we requested it due to the swap out.
			return resp, "Pending", nil
		}

		return resp, "Updated", nil
	}
}

func expandRoleDefinitionPermissions(input []interface{}) []authorization.Permission {
	output := make([]authorization.Permission, 0)
	if len(input) == 0 {
		return output
	}

	for _, v := range input {
		if v == nil {
			continue
		}

		raw := v.(map[string]interface{})
		permission := authorization.Permission{}

		actionsOutput := make([]string, 0)
		actions := raw["actions"].([]interface{})
		for _, a := range actions {
			if a == nil {
				continue
			}
			actionsOutput = append(actionsOutput, a.(string))
		}
		permission.Actions = &actionsOutput

		dataActionsOutput := make([]string, 0)
		dataActions := raw["data_actions"].(*pluginsdk.Set)
		for _, a := range dataActions.List() {
			if a == nil {
				continue
			}
			dataActionsOutput = append(dataActionsOutput, a.(string))
		}
		permission.DataActions = &dataActionsOutput

		notActionsOutput := make([]string, 0)
		notActions := raw["not_actions"].([]interface{})
		for _, a := range notActions {
			if a == nil {
				continue
			}
			notActionsOutput = append(notActionsOutput, a.(string))
		}
		permission.NotActions = &notActionsOutput

		notDataActionsOutput := make([]string, 0)
		notDataActions := raw["not_data_actions"].(*pluginsdk.Set)
		for _, a := range notDataActions.List() {
			if a == nil {
				continue
			}
			notDataActionsOutput = append(notDataActionsOutput, a.(string))
		}
		permission.NotDataActions = &notDataActionsOutput

		output = append(output, permission)
	}

	return output
}

func expandRoleDefinitionAssignableScopes(d *pluginsdk.ResourceData) []string {
	scopes := make([]string, 0)

	assignableScopes := d.Get("assignable_scopes").([]interface{})
	if len(assignableScopes) == 0 {
		assignedScope := d.Get("scope").(string)
		scopes = append(scopes, assignedScope)
	} else {
		for _, scope := range assignableScopes {
			scopes = append(scopes, scope.(string))
		}
	}

	return scopes
}

func flattenRoleDefinitionPermissions(input *[]authorization.Permission) []interface{} {
	permissions := make([]interface{}, 0)
	if input == nil {
		return permissions
	}

	for _, permission := range *input {
		permissions = append(permissions, map[string]interface{}{
			"actions":          utils.FlattenStringSlice(permission.Actions),
			"data_actions":     pluginsdk.NewSet(pluginsdk.HashString, utils.FlattenStringSlice(permission.DataActions)),
			"not_actions":      utils.FlattenStringSlice(permission.NotActions),
			"not_data_actions": pluginsdk.NewSet(pluginsdk.HashString, utils.FlattenStringSlice(permission.NotDataActions)),
		})
	}

	return permissions
}

func flattenRoleDefinitionAssignableScopes(input *[]string) []interface{} {
	scopes := make([]interface{}, 0)
	if input == nil {
		return scopes
	}

	for _, scope := range *input {
		scopes = append(scopes, scope)
	}

	return scopes
}

func roleDefinitionUpdateStateRefreshFunc(ctx context.Context, client *authorization.RoleDefinitionsClient, roleDefinitionId string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.GetByID(ctx, roleDefinitionId)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return resp, "NotFound", err
			}
			return resp, "Error", err
		}
		return "OK", "OK", nil
	}
}

func roleDefinitionDeleteStateRefreshFunc(ctx context.Context, client *authorization.RoleDefinitionsClient, roleDefinitionId string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.GetByID(ctx, roleDefinitionId)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return resp, "NotFound", nil
			}
			return nil, "Error", err
		}
		return "Pending", "Pending", nil
	}
}
