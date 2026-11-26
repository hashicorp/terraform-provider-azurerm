// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package policy_definition

import (
	"fmt"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePolicyDefinitionDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Policy.DefinitionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PolicyDefinitionID(d.Id())
	if err != nil {
		return err
	}

	managementGroupName := ""
	switch scopeId := id.PolicyScopeId.(type) { //nolint:gocritic
	case parse.ScopeAtManagementGroup:
		managementGroupName = scopeId.ManagementGroupName
	}

	var resp autorest.Response
	if managementGroupName == "" {
		resp, err = client.Delete(ctx, id.Name)
	} else {
		resp, err = client.DeleteAtManagementGroup(ctx, id.Name, managementGroupName)
	}

	if err != nil {
		if response.WasNotFound(resp.Response) {
			return nil
		}

		return fmt.Errorf("deleting Policy Definition %q: %+v", id.Name, err)
	}

	return nil
}
