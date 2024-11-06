// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

var _ plancheck.PlanCheck = isNotResourceAction{}

type isNotResourceAction struct {
	resourceAddress string
	actionType      plancheck.ResourceActionType
}

// CheckPlan implements the plan check logic.
func (e isNotResourceAction) CheckPlan(ctx context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
	foundResource := false

	for _, rc := range req.Plan.ResourceChanges {
		if e.resourceAddress != rc.Address {
			continue
		}

		switch e.actionType {
		case plancheck.ResourceActionNoop:
			if rc.Change.Actions.NoOp() {
				resp.Error = fmt.Errorf("'%s' - expected action to not be %s", rc.Address, e.actionType)
				return
			}
		case plancheck.ResourceActionCreate:
			if rc.Change.Actions.Create() {
				resp.Error = fmt.Errorf("'%s' - expected action to not be %s", rc.Address, e.actionType)
				return
			}
		case plancheck.ResourceActionRead:
			if rc.Change.Actions.Read() {
				resp.Error = fmt.Errorf("'%s' - expected action to not be %s", rc.Address, e.actionType)
				return
			}
		case plancheck.ResourceActionUpdate:
			if rc.Change.Actions.Update() {
				resp.Error = fmt.Errorf("'%s' - expected action to not be %s", rc.Address, e.actionType)
				return
			}
		case plancheck.ResourceActionDestroy:
			if rc.Change.Actions.Delete() {
				resp.Error = fmt.Errorf("'%s' - expected action to not be %s", rc.Address, e.actionType)
				return
			}
		case plancheck.ResourceActionDestroyBeforeCreate:
			if rc.Change.Actions.DestroyBeforeCreate() {
				resp.Error = fmt.Errorf("'%s' - expected action to not be %s", rc.Address, e.actionType)
				return
			}
		case plancheck.ResourceActionCreateBeforeDestroy:
			if rc.Change.Actions.CreateBeforeDestroy() {
				resp.Error = fmt.Errorf("'%s' - expected action to not be %s", rc.Address, e.actionType)
				return
			}
		case plancheck.ResourceActionReplace:
			if rc.Change.Actions.Replace() {
				resp.Error = fmt.Errorf("'%s' - expected action to not be %s, path: %v", rc.Address, e.actionType, rc.Change.ReplacePaths)
				return
			}
		default:
			resp.Error = fmt.Errorf("%s - unexpected ResourceActionType: %s", rc.Address, e.actionType)
			return
		}

		foundResource = true
		break
	}

	if !foundResource {
		resp.Error = fmt.Errorf("%s - Resource not found in plan ResourceChanges", e.resourceAddress)
		return
	}
}

// IsNotResourceAction returns a plan check that asserts that a given resource will not have a specific resource change type in the plan.
// Valid actionType are an enum of type plancheck.ResourceActionType, examples: NoOp, DestroyBeforeCreate, Update (in-place), etc.
func IsNotResourceAction(resourceAddress string, actionType plancheck.ResourceActionType) plancheck.PlanCheck {
	return isNotResourceAction{
		resourceAddress: resourceAddress,
		actionType:      actionType,
	}
}
