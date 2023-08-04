// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package blueprints

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/blueprints/2018-11-01-preview/assignment"
	"github.com/hashicorp/go-azure-sdk/resource-manager/blueprints/2018-11-01-preview/publishedblueprint"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceBlueprintAssignment() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceBlueprintAssignmentCreateUpdate,
		Update: resourceBlueprintAssignmentCreateUpdate,
		Read:   resourceBlueprintAssignmentRead,
		Delete: resourceBlueprintAssignmentDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := assignment.ParseScopedBlueprintAssignmentID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"target_subscription_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"location": commonschema.Location(),

			"identity": commonschema.SystemOrUserAssignedIdentityRequired(),

			"version_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: publishedblueprint.ValidateScopedVersionID,
			},

			"parameter_values": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				StateFunc:        normalizeAssignmentParameterValuesJSON,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
			},

			"resource_groups": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				StateFunc:        normalizeAssignmentResourceGroupValuesJSON,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
			},

			"lock_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(assignment.AssignmentLockModeNone),
				ValidateFunc: validation.StringInSlice([]string{
					string(assignment.AssignmentLockModeNone),
					string(assignment.AssignmentLockModeAllResourcesReadOnly),
					string(assignment.AssignmentLockModeAllResourcesDoNotDelete),
				}, false),
				// The first character of value returned by the service is always in lower case.
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"lock_exclude_principals": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 5,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.IsUUID,
				},
			},

			"lock_exclude_actions": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 200,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"display_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"blueprint_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceBlueprintAssignmentCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Blueprints.AssignmentsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := assignment.NewScopedBlueprintAssignmentID(d.Get("target_subscription_id").(string), d.Get("name").(string))
	blueprintId := d.Get("version_id").(string)

	if d.IsNewResource() {
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
		assignmentLockSettings.Mode = pointer.To(assignment.AssignmentLockMode(lockMode))
		if lockMode != "None" {
			excludedPrincipalsRaw := d.Get("lock_exclude_principals").([]interface{})
			if len(excludedPrincipalsRaw) != 0 {
				assignmentLockSettings.ExcludedPrincipals = utils.ExpandStringSlice(excludedPrincipalsRaw)
			}

			excludedActionsRaw := d.Get("lock_exclude_actions").([]interface{})
			if len(excludedActionsRaw) != 0 {
				assignmentLockSettings.ExcludedActions = utils.ExpandStringSlice(excludedActionsRaw)
			}
		}
		payload.Properties.Locks = assignmentLockSettings
	}

	i, err := identity.ExpandSystemOrUserAssignedMap(d.Get("identity").([]interface{}))
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

	d.SetId(id.ID())

	return resourceBlueprintAssignmentRead(d, meta)
}

func resourceBlueprintAssignmentRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

		return fmt.Errorf("Read failed for Blueprint Assignment (%q): %+v", id.String(), err)
	}

	d.Set("name", id.BlueprintAssignmentName)
	if model := resp.Model; model != nil {
		p := model.Properties

		d.Set("location", azure.NormalizeLocation(model.Location))
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
			d.Set("lock_mode", string(pointer.From(locks.Mode)))
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

func resourceBlueprintAssignmentDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Blueprints.AssignmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := assignment.ParseScopedBlueprintAssignmentID(d.Id())
	if err != nil {
		return err
	}

	// We use none here to align the previous behaviour of the blueprint resource
	// TODO: we could add a features flag for the blueprint to empower terraform when deleting the blueprint to delete all the generated resources as well
	if _, err := client.Delete(ctx, *id, assignment.DeleteOperationOptions{}); err != nil {
		return fmt.Errorf("failed to delete Blueprint Assignment %q from scope %q: %+v", id.BlueprintAssignmentName, id.ResourceScope, err)
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{
			string(assignment.AssignmentProvisioningStateWaiting),
			string(assignment.AssignmentProvisioningStateValidating),
			string(assignment.AssignmentProvisioningStateLocking),
			string(assignment.AssignmentProvisioningStateDeleting),
			string(assignment.AssignmentProvisioningStateFailed),
		},
		Target:  []string{"NotFound"},
		Refresh: blueprintAssignmentDeleteStateRefreshFunc(ctx, client, *id),
		Timeout: time.Until(deadline),
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for Blueprint Assignment %q: %+v", id.String(), err)
	}

	return nil
}
