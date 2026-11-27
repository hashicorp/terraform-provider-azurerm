// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_role_assignment

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	synapseClient2 "github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/synapse/2020-08-01-preview/synapse"
)

func resourceSynapseRoleAssignmentRead(d *pluginsdk.ResourceData, meta any) error {
	synapseClient := synapseClient2.From(meta)
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	env := account.From(meta).Environment
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
		if response.WasNotFound(resp.Response.Response) {
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

	d.Set("principal_type", pointer.From(resp.PrincipalType))

	synapseWorkspaceId := ""
	synapseSparkPoolId := ""
	if _, err := workspaces.ParseWorkspaceIDInsensitively(id.Scope); err == nil {
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

func getRoleIdByName(ctx context.Context, client *synapse.RoleDefinitionsClient, scope, roleName string) (*uuid.UUID, error) {
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
