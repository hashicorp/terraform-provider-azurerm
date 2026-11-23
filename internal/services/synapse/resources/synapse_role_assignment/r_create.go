// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_role_assignment

import (
	"fmt"

	frsUUID "github.com/gofrs/uuid"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/resources/synapse_role_assignment/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	accesscontrol "github.com/jackofallops/kermit/sdk/synapse/2020-08-01-preview/synapse"
)

func resourceSynapseRoleAssignmentCreate(d *pluginsdk.ResourceData, meta any) error {
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
		if !response.WasNotFound(listResp.Response.Response) {
			return fmt.Errorf("checking for presence of existing Synapse Role Assignment (workspace %q): %+v", workspaceName, err)
		}
	}
	// TODO: unpick this/refactor to use ID Formatters
	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		if listResp.Value != nil && len(*listResp.Value) != 0 {
			existing := (*listResp.Value)[0]
			if !response.WasNotFound(existing.Response.Response) {
				resourceId := parse.NewRoleAssignmentId(synapseScope, *existing.ID).ID()
				return tf.ImportAsExistsError("azurerm_synapse_role_assignment", resourceId)
			}
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
		Scope:       pointer.To(scope),
	}

	if v, ok := d.GetOk("principal_type"); ok {
		roleAssignment.PrincipalType = pointer.To(v.(string))
	}

	resp, err := client.CreateRoleAssignment(ctx, roleAssignment, uuid)
	if err != nil {
		return fmt.Errorf("creating Synapse RoleAssignment %q: %+v", roleName, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Synapse RoleAssignment %q", roleName)
	}

	d.SetId(parse.NewRoleAssignmentId(synapseScope, *resp.ID).ID())
	return resourceSynapseRoleAssignmentRead(d, meta)
}
