// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plancheck

import (
	"context"
	"fmt"
)

var _ PlanCheck = expectResourceAction{}

type expectResourceAction struct {
	resourceAddress string
	actionType      ResourceActionType
}

// CheckPlan implements the plan check logic.
func (e expectResourceAction) CheckPlan(ctx context.Context, req CheckPlanRequest, resp *CheckPlanResponse) {
	foundResource := false

	for _, rc := range req.Plan.ResourceChanges {
		if e.resourceAddress != rc.Address {
			continue
		}

		switch e.actionType {
		case ResourceActionNoop:
			if !rc.Change.Actions.NoOp() {
				resp.Error = fmt.Errorf("'%s' - expected %s, got action(s): %v", rc.Address, e.actionType, rc.Change.Actions)
				return
			}
		case ResourceActionCreate:
			if !rc.Change.Actions.Create() {
				resp.Error = fmt.Errorf("'%s' - expected %s, got action(s): %v", rc.Address, e.actionType, rc.Change.Actions)
				return
			}
		case ResourceActionRead:
			if !rc.Change.Actions.Read() {
				resp.Error = fmt.Errorf("'%s' - expected %s, got action(s): %v", rc.Address, e.actionType, rc.Change.Actions)
				return
			}
		case ResourceActionUpdate:
			if !rc.Change.Actions.Update() {
				resp.Error = fmt.Errorf("'%s' - expected %s, got action(s): %v", rc.Address, e.actionType, rc.Change.Actions)
				return
			}
		case ResourceActionDestroy:
			if !rc.Change.Actions.Delete() {
				resp.Error = fmt.Errorf("'%s' - expected %s, got action(s): %v", rc.Address, e.actionType, rc.Change.Actions)
				return
			}
		case ResourceActionDestroyBeforeCreate:
			if !rc.Change.Actions.DestroyBeforeCreate() {
				resp.Error = fmt.Errorf("'%s' - expected %s, got action(s): %v", rc.Address, e.actionType, rc.Change.Actions)
				return
			}
		case ResourceActionCreateBeforeDestroy:
			if !rc.Change.Actions.CreateBeforeDestroy() {
				resp.Error = fmt.Errorf("'%s' - expected %s, got action(s): %v", rc.Address, e.actionType, rc.Change.Actions)
				return
			}
		case ResourceActionReplace:
			if !rc.Change.Actions.Replace() {
				resp.Error = fmt.Errorf("%s - expected %s, got action(s): %v", rc.Address, e.actionType, rc.Change.Actions)
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

// ExpectResourceAction returns a plan check that asserts that a given resource will have a specific resource change type in the plan.
// Valid actionType are an enum of type plancheck.ResourceActionType, examples: NoOp, DestroyBeforeCreate, Update (in-place), etc.
func ExpectResourceAction(resourceAddress string, actionType ResourceActionType) PlanCheck {
	return expectResourceAction{
		resourceAddress: resourceAddress,
		actionType:      actionType,
	}
}
