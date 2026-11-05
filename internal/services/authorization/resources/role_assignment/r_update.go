// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package role_assignment

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2022-04-01/roleassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmRoleAssignmentUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Authorization.ScopedRoleAssignmentsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ScopedRoleAssignmentID(d.Id())
	if err != nil {
		return err
	}

	options := roleassignments.DefaultGetOperationOptions()
	if id.TenantId != "" {
		options.TenantId = pointer.To(id.TenantId)
	}

	existing, err := client.Get(ctx, id.ScopedId, options)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if existing.Model == nil || existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `model.properties` was nil", *id)
	}

	props := *existing.Model.Properties

	if d.HasChange("description") {
		props.Description = nil
		if v := d.Get("description").(string); v != "" {
			props.Description = &v
		}
	}

	// Order matters here that "condition_version" shall be handled prior to "condition".
	if d.HasChange("condition_version") {
		props.ConditionVersion = nil
		if v := d.Get("condition_version").(string); v != "" {
			props.ConditionVersion = &v
		}
	}
	if d.HasChange("condition") {
		props.Condition = nil
		if v := d.Get("condition").(string); v != "" {
			props.Condition = &v
		}

		// Implicitly setting the condition_version in case it is not specified in config (as how Create() has been implemented).
		if d.GetRawConfig().AsValueMap()["condition_version"].IsNull() {
			if props.Condition == nil {
				props.ConditionVersion = nil
			} else {
				props.ConditionVersion = pointer.To("2.0")
			}
		}
	}

	params := roleassignments.RoleAssignmentCreateParameters{
		Properties: props,
	}

	if _, err := client.Create(ctx, id.ScopedId, params); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceArmRoleAssignmentRead(d, meta)
}
