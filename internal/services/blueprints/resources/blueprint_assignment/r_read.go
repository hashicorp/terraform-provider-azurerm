// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package blueprint_assignment

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/blueprints/2018-11-01-preview/assignment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceBlueprintAssignmentRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Blueprints.AssignmentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := assignment.ParseScopedBlueprintAssignmentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] the Blueprint Assignment %q does not exist - removing from state", id.String())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("read failed for Blueprint Assignment (%q): %+v", id.String(), err)
	}

	d.Set("name", id.BlueprintAssignmentName)
	if model := resp.Model; model != nil {
		p := model.Properties

		d.Set("location", location.Normalize(model.Location))
		d.Set("target_subscription_id", pointer.From(p.Scope))
		d.Set("version_id", pointer.From(p.BlueprintId))
		d.Set("display_name", pointer.From(p.DisplayName))
		d.Set("description", pointer.From(p.Description))

		if p.Parameters != nil {
			params, err := flattenArmBlueprintAssignmentParameters(p.Parameters)
			if err != nil {
				return err
			}
			d.Set("parameter_values", params)
		}

		if p.ResourceGroups != nil {
			resourceGroups, err := flattenArmBlueprintAssignmentResourceGroups(p.ResourceGroups)
			if err != nil {
				return err
			}
			d.Set("resource_groups", resourceGroups)
		}

		// Locks
		if locks := p.Locks; locks != nil {
			d.Set("lock_mode", pointer.FromEnum(locks.Mode))
			if locks.ExcludedPrincipals != nil {
				d.Set("lock_exclude_principals", locks.ExcludedPrincipals)
			}
			if locks.ExcludedActions != nil {
				d.Set("lock_exclude_actions", locks.ExcludedActions)
			}
		}

		i, err := identity.FlattenSystemOrUserAssignedMap(&model.Identity)
		if err != nil {
			return err
		}
		if err := d.Set("identity", i); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}
	}

	return nil
}
