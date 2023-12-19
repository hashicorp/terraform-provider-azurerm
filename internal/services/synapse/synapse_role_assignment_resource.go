// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package synapse

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	frsUUID "github.com/gofrs/uuid"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	accesscontrol "github.com/tombuildsstuff/kermit/sdk/synapse/2020-08-01-preview/synapse"
)

func resourceSynapseRoleAssignment() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSynapseRoleAssignmentCreate,
		Read:   resourceSynapseRoleAssignmentRead,
		Delete: resourceSynapseRoleAssignmentDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.RoleAssignmentV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.RoleAssignmentID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"synapse_workspace_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"synapse_workspace_id", "synapse_spark_pool_id"},
				ValidateFunc: validate.WorkspaceID,
			},

			"synapse_spark_pool_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"synapse_workspace_id", "synapse_spark_pool_id"},
				ValidateFunc: validate.SparkPoolID,
			},

			"principal_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"principal_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"User",
					"Group",
					"ServicePrincipal",
				}, false),
			},

			"role_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(_, old, new string, d *pluginsdk.ResourceData) bool {
					return migration.MigrateToNewRole(old) == migration.MigrateToNewRole(new)
				},
				ValidateFunc: validation.StringInSlice([]string{
					"Apache Spark Administrator",
					"Synapse Administrator",
					"Synapse Artifact Publisher",
					"Synapse Artifact User",
					"Synapse Compute Operator",
					"Synapse Contributor",
					"Synapse Credential User",
					"Synapse Linked Data Manager",
					"Synapse Monitoring Operator",
					"Synapse SQL Administrator",
					"Synapse User",
				}, false),
			},
		},
	}
}

func resourceSynapseRoleAssignmentCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	env := meta.(*clients.Client).Account.Environment
	synapseDomainSuffix, ok := env.Synapse.DomainSuffix()
	if !ok {
		return fmt.Errorf("could not determine the domain suffix for synapse in environment %q: %+v", env.Name, env.Storage)
	}

	synapseScope := ""
	if v, ok := d.GetOk("synapse_workspace_id"); ok {
		synapseScope = v.(string)
	} else if v, ok := d.GetOk("synapse_spark_pool_id"); ok {
		synapseScope = v.(string)
	}

	workspaceName, scope, err := parse.SynapseScope(synapseScope)
	if err != nil {
		return err
	}

	client, err := synapseClient.RoleAssignmentsClient(workspaceName, *synapseDomainSuffix)
	if err != nil {
		return err
	}
	roleDefinitionsClient, err := synapseClient.RoleDefinitionsClient(workspaceName, *synapseDomainSuffix)
	if err != nil {
		return err
	}

	roleName := migration.MigrateToNewRole(d.Get("role_name").(string))
	roleId, err := getRoleIdByName(ctx, roleDefinitionsClient, scope, roleName)
	if err != nil {
		return err
	}

	// check exist
	principalId := d.Get("principal_id").(string)
	listResp, err := client.ListRoleAssignments(ctx, roleId.String(), principalId, scope, "")
	if err != nil {
		if !utils.ResponseWasNotFound(listResp.Response) {
			return fmt.Errorf("checking for presence of existing Synapse Role Assignment (workspace %q): %+v", workspaceName, err)
		}
	}
	// TODO: unpick this/refactor to use ID Formatters
	if listResp.Value != nil && len(*listResp.Value) != 0 {
		existing := (*listResp.Value)[0]
		if !utils.ResponseWasNotFound(existing.Response) {
			resourceId := parse.NewRoleAssignmentId(synapseScope, *existing.ID).ID()
			return tf.ImportAsExistsError("azurerm_synapse_role_assignment", resourceId)
		}
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return fmt.Errorf("generating UUID for Synapse Role Assignment: %+v", err)
	}

	principalID, err := frsUUID.FromString(principalId)
	if err != nil {
		return err
	}

	// create
	roleAssignment := accesscontrol.RoleAssignmentRequest{
		RoleID:      roleId,
		PrincipalID: &principalID,
		Scope:       utils.String(scope),
	}

	if v, ok := d.GetOk("principal_type"); ok {
		principalType := v.(string)
		roleAssignment.PrincipalType = &principalType
	}

	resp, err := client.CreateRoleAssignment(ctx, roleAssignment, uuid)
	if err != nil {
		return fmt.Errorf("creating Synapse RoleAssignment %q: %+v", roleName, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Synapse RoleAssignment %q", roleName)
	}

	resourceId := parse.NewRoleAssignmentId(synapseScope, *resp.ID).ID()
	d.SetId(resourceId)
	return resourceSynapseRoleAssignmentRead(d, meta)
}

func resourceSynapseRoleAssignmentRead(d *pluginsdk.ResourceData, meta interface{}) error {
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	env := meta.(*clients.Client).Account.Environment
	synapseDomainSuffix, ok := env.Synapse.DomainSuffix()
	if !ok {
		return fmt.Errorf("could not determine the domain suffix for synapse in environment %q: %+v", env.Name, env.Storage)
	}

	id, err := parse.RoleAssignmentID(d.Id())
	if err != nil {
		return err
	}

	workspaceName, _, err := parse.SynapseScope(id.Scope)
	if err != nil {
		return err
	}

	client, err := synapseClient.RoleAssignmentsClient(workspaceName, *synapseDomainSuffix)
	if err != nil {
		return err
	}
	roleDefinitionsClient, err := synapseClient.RoleDefinitionsClient(workspaceName, *synapseDomainSuffix)
	if err != nil {
		return err
	}

	resp, err := client.GetRoleAssignmentByID(ctx, id.DataPlaneAssignmentId)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] synapse role assignment %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Synapse RoleAssignment (Resource Group %q): %+v", workspaceName, err)
	}

	principalID := ""
	if resp.PrincipalID != nil {
		principalID = resp.PrincipalID.String()
	}
	d.Set("principal_id", principalID)

	principalType := ""
	if resp.PrincipalType != nil {
		principalType = *resp.PrincipalType
	}
	d.Set("principal_type", principalType)

	synapseWorkspaceId := ""
	synapseSparkPoolId := ""
	if _, err := parse.WorkspaceIDInsensitively(id.Scope); err == nil {
		synapseWorkspaceId = id.Scope
	} else if _, err := parse.SparkPoolIDInsensitively(id.Scope); err == nil {
		synapseSparkPoolId = id.Scope
	}

	d.Set("synapse_workspace_id", synapseWorkspaceId)
	d.Set("synapse_spark_pool_id", synapseSparkPoolId)

	if resp.RoleDefinitionID != nil {
		role, err := roleDefinitionsClient.GetRoleDefinitionByID(ctx, resp.RoleDefinitionID.String())
		if err != nil {
			return fmt.Errorf("retrieving role definition by ID %q: %+v", resp.RoleDefinitionID.String(), err)
		}
		d.Set("role_name", role.Name)
	}
	return nil
}

func resourceSynapseRoleAssignmentDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	env := meta.(*clients.Client).Account.Environment
	synapseDomainSuffix, ok := env.Synapse.DomainSuffix()
	if !ok {
		return fmt.Errorf("could not determine the domain suffix for synapse in environment %q: %+v", env.Name, env.Storage)
	}

	id, err := parse.RoleAssignmentID(d.Id())
	if err != nil {
		return err
	}

	workspaceName, scope, err := parse.SynapseScope(id.Scope)
	if err != nil {
		return err
	}

	client, err := synapseClient.RoleAssignmentsClient(workspaceName, *synapseDomainSuffix)
	if err != nil {
		return err
	}
	if _, err := client.DeleteRoleAssignmentByID(ctx, id.DataPlaneAssignmentId, scope); err != nil {
		return fmt.Errorf("deleting Synapse RoleAssignment %q (workspace %q): %+v", id, workspaceName, err)
	}

	return nil
}

func getRoleIdByName(ctx context.Context, client *accesscontrol.RoleDefinitionsClient, scope, roleName string) (*frsUUID.UUID, error) {
	resp, err := client.ListRoleDefinitions(ctx, nil, scope)
	if err != nil {
		return nil, fmt.Errorf("listing synapse role definitions %+v", err)
	}

	var availableRoleName []string
	if resp.Value != nil {
		for _, role := range *resp.Value {
			if role.Name != nil {
				if *role.Name == roleName && role.ID != nil {
					return role.ID, nil
				}
				availableRoleName = append(availableRoleName, *role.Name)
			}
		}
	}

	return nil, fmt.Errorf("role name %q invalid for scope %q. Available role names are %q", roleName, scope, strings.Join(availableRoleName, ","))
}
