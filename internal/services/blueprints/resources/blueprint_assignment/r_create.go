// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package blueprint_assignment

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/blueprints/2018-11-01-preview/assignment"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceBlueprintAssignmentCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Blueprints.AssignmentsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := assignment.NewScopedBlueprintAssignmentID(d.Get("target_subscription_id").(string), d.Get("name").(string))
	blueprintId := d.Get("version_id").(string)

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			resp, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("checking for an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(resp.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_blueprint_assignment", id.ID())
			}
		}
	}

	payload := assignment.Assignment{
		Properties: assignment.AssignmentProperties{
			BlueprintId: pointer.To(blueprintId), // This is mislabeled - The ID is that of the Published Version, not just the Blueprint
			Scope:       pointer.To(id.ResourceScope),
		},
		Location: location.Normalize(d.Get("location").(string)),
	}

	if lockModeRaw, ok := d.GetOk("lock_mode"); ok {
		assignmentLockSettings := &assignment.AssignmentLockSettings{}
		lockMode := lockModeRaw.(string)
		assignmentLockSettings.Mode = pointer.ToEnum[assignment.AssignmentLockMode](lockMode)
		if lockMode != "None" {
			excludedPrincipalsRaw := d.Get("lock_exclude_principals").([]any)
			if len(excludedPrincipalsRaw) != 0 {
				assignmentLockSettings.ExcludedPrincipals = pluginsdk.ExpandStringSlice(excludedPrincipalsRaw)
			}

			excludedActionsRaw := d.Get("lock_exclude_actions").([]any)
			if len(excludedActionsRaw) != 0 {
				assignmentLockSettings.ExcludedActions = pluginsdk.ExpandStringSlice(excludedActionsRaw)
			}
		}
		payload.Properties.Locks = assignmentLockSettings
	}

	i, err := identity.ExpandSystemOrUserAssignedMap(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}
	payload.Identity = *i

	if paramValuesRaw := d.Get("parameter_values"); paramValuesRaw != "" {
		payload.Properties.Parameters = expandArmBlueprintAssignmentParameters(paramValuesRaw.(string))
	} else {
		payload.Properties.Parameters = expandArmBlueprintAssignmentParameters("{}")
	}

	if resourceGroupsRaw := d.Get("resource_groups"); resourceGroupsRaw != "" {
		payload.Properties.ResourceGroups = expandArmBlueprintAssignmentResourceGroups(resourceGroupsRaw.(string))
	} else {
		payload.Properties.ResourceGroups = expandArmBlueprintAssignmentResourceGroups("{}")
	}

	if _, err = client.CreateOrUpdate(ctx, id, payload); err != nil {
		return err
	}

	d.SetId(id.ID())

	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{
			string(assignment.AssignmentProvisioningStateWaiting),
			string(assignment.AssignmentProvisioningStateValidating),
			string(assignment.AssignmentProvisioningStateCreating),
			string(assignment.AssignmentProvisioningStateDeploying),
			string(assignment.AssignmentProvisioningStateLocking),
		},
		Target:  []string{string(assignment.AssignmentProvisioningStateSucceeded)},
		Refresh: blueprintAssignmentCreateStateRefreshFunc(ctx, client, id),
		Timeout: time.Until(deadline),
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("failed waiting for Blueprint Assignment %s: %+v", id.String(), err)
	}

	return resourceBlueprintAssignmentRead(d, meta)
}
