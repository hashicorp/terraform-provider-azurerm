// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package synapse

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/data-plane/synapse/2020-08-01-preview/roleassignments"
	"github.com/hashicorp/go-azure-sdk/data-plane/synapse/2020-08-01-preview/synapseroledefinitions"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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

	endpoint := fmt.Sprintf("https://%s.%s", workspaceName, *synapseDomainSuffix)

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
	listResp, err := client.ListRoleAssignments(ctx, roleassignments.ListRoleAssignmentsOperationOptions{
		RoleId:      pointer.To(roleId),
		PrincipalId: pointer.To(principalId),
		Scope:       pointer.To(scope),
	})
	if err != nil {
		if !response.WasNotFound(listResp.HttpResponse) {
			return fmt.Errorf("checking for presence of existing Synapse Role Assignment (workspace %q): %+v", workspaceName, err)
		}
	}
	// TODO: unpick this/refactor to use ID Formatters
	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		if listResp.Model != nil && listResp.Model.Value != nil && len(*listResp.Model.Value) != 0 {
			existing := (*listResp.Model.Value)[0]
			if existing.Id != nil && *existing.Id != "" {
				resourceId := parse.NewRoleAssignmentId(synapseScope, *existing.Id).ID()
				return tf.ImportAsExistsError("azurerm_synapse_role_assignment", resourceId)
			}
		}
	}

	assignmentUuid, err := uuid.GenerateUUID()
	if err != nil {
		return fmt.Errorf("generating UUID for Synapse Role Assignment: %+v", err)
	}

	// create
	roleAssignment := roleassignments.RoleAssignmentRequest{
		RoleId:      roleId,
		PrincipalId: principalId,
		Scope:       scope,
	}

	if v, ok := d.GetOk("principal_type"); ok {
		roleAssignment.PrincipalType = pointer.To(v.(string))
	}

	assignmentId := roleassignments.NewRoleAssignmentIdID(endpoint, assignmentUuid)
	resp, err := client.CreateRoleAssignment(ctx, assignmentId, roleAssignment)
	if err != nil {
		return fmt.Errorf("creating Synapse RoleAssignment %q: %+v", roleName, err)
	}

	if resp.Model == nil || resp.Model.Id == nil || *resp.Model.Id == "" {
		return fmt.Errorf("empty or nil ID returned for Synapse RoleAssignment %q", roleName)
	}

	resourceId := parse.NewRoleAssignmentId(synapseScope, *resp.Model.Id).ID()
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

	endpoint := fmt.Sprintf("https://%s.%s", workspaceName, *synapseDomainSuffix)

	client, err := synapseClient.RoleAssignmentsClient(workspaceName, *synapseDomainSuffix)
	if err != nil {
		return err
	}
	roleDefinitionsClient, err := synapseClient.RoleDefinitionsClient(workspaceName, *synapseDomainSuffix)
	if err != nil {
		return err
	}

	resp, err := client.GetRoleAssignmentById(ctx, roleassignments.NewRoleAssignmentIdID(endpoint, id.DataPlaneAssignmentId))
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] synapse role assignment %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Synapse RoleAssignment (Resource Group %q): %+v", workspaceName, err)
	}

	synapseWorkspaceId := ""
	synapseSparkPoolId := ""
	if _, err := parse.WorkspaceIDInsensitively(id.Scope); err == nil {
		synapseWorkspaceId = id.Scope
	} else if _, err := parse.SparkPoolIDInsensitively(id.Scope); err == nil {
		synapseSparkPoolId = id.Scope
	}

	d.Set("synapse_workspace_id", synapseWorkspaceId)
	d.Set("synapse_spark_pool_id", synapseSparkPoolId)

	if model := resp.Model; model != nil {
		d.Set("principal_id", pointer.From(model.PrincipalId))
		d.Set("principal_type", pointer.From(model.PrincipalType))

		if model.RoleDefinitionId != nil {
			roleDefinitionId := *model.RoleDefinitionId
			role, err := roleDefinitionsClient.RoleDefinitionsGetRoleDefinitionById(ctx, synapseroledefinitions.NewRoleDefinitionID(endpoint, roleDefinitionId))
			if err != nil {
				return fmt.Errorf("retrieving role definition by ID %q: %+v", roleDefinitionId, err)
			}
			if role.Model != nil {
				d.Set("role_name", pointer.From(role.Model.Name))
			}
		}
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

	endpoint := fmt.Sprintf("https://%s.%s", workspaceName, *synapseDomainSuffix)

	client, err := synapseClient.RoleAssignmentsClient(workspaceName, *synapseDomainSuffix)
	if err != nil {
		return err
	}
	if _, err := client.DeleteRoleAssignmentById(ctx, roleassignments.NewRoleAssignmentIdID(endpoint, id.DataPlaneAssignmentId), roleassignments.DeleteRoleAssignmentByIdOperationOptions{
		Scope: pointer.To(scope),
	}); err != nil {
		return fmt.Errorf("deleting Synapse RoleAssignment %q (workspace %q): %+v", id, workspaceName, err)
	}

	return nil
}

func getRoleIdByName(ctx context.Context, client *synapseroledefinitions.SynapseRoleDefinitionsClient, scope, roleName string) (string, error) {
	resp, err := client.RoleDefinitionsListRoleDefinitions(ctx, synapseroledefinitions.RoleDefinitionsListRoleDefinitionsOperationOptions{
		Scope: pointer.To(scope),
	})
	if err != nil {
		return "", fmt.Errorf("listing synapse role definitions %+v", err)
	}

	var availableRoleName []string
	if resp.Model != nil {
		for _, role := range *resp.Model {
			if role.Name != nil {
				if *role.Name == roleName && role.Id != nil {
					return *role.Id, nil
				}
				availableRoleName = append(availableRoleName, *role.Name)
			}
		}
	}

	return "", fmt.Errorf("role name %q invalid for scope %q. Available role names are %q", roleName, scope, strings.Join(availableRoleName, ","))
}
