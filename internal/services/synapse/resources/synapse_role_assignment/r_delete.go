// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_role_assignment

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	synapseClient2 "github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseRoleAssignmentDelete(d *pluginsdk.ResourceData, meta any) error {
	synapseClient := synapseClient2.From(meta)
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
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
