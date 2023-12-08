// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package blueprints

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/blueprints/2018-11-01-preview/assignment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func blueprintAssignmentCreateStateRefreshFunc(ctx context.Context, client *assignment.AssignmentClient, id assignment.ScopedBlueprintAssignmentId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, id)
		if err != nil {
			return nil, "", fmt.Errorf("unable to retrieve Blueprint Assignment %s: %+v", id.String(), err)
		}
		if resp.Model == nil || resp.Model.Properties.ProvisioningState == nil {
			return resp, "nil", fmt.Errorf("Blueprint Assignment Model or ProvisioningState is nil")
		}
		state := *resp.Model.Properties.ProvisioningState

		if state == assignment.AssignmentProvisioningStateFailed {
			return resp, string(state), fmt.Errorf("Blueprint Assignment provisioning entered a Failed state.")
		}

		return resp, string(state), nil
	}
}

func blueprintAssignmentDeleteStateRefreshFunc(ctx context.Context, client *assignment.AssignmentClient, id assignment.ScopedBlueprintAssignmentId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, id)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return resp, "NotFound", nil
			} else {
				return nil, "", fmt.Errorf("unable to retrieve Blueprint Assignment %s: %+v", id.String(), err)
			}
		}

		if resp.Model == nil || resp.Model.Properties.ProvisioningState == nil {
			return resp, "nil", fmt.Errorf("Blueprint Assignment Model or ProvisioningState is nil")
		}
		return resp, string(*resp.Model.Properties.ProvisioningState), nil
	}
}

func normalizeAssignmentParameterValuesJSON(jsonString interface{}) string {
	if jsonString == nil || jsonString == "" {
		return ""
	}

	var values map[string]*assignment.ParameterValue
	if err := json.Unmarshal([]byte(jsonString.(string)), &values); err != nil {
		return fmt.Sprintf("unable to parse JSON: %+v", err)
	}

	b, _ := json.Marshal(values)
	return string(b)
}

func normalizeAssignmentResourceGroupValuesJSON(jsonString interface{}) string {
	if jsonString == nil || jsonString == "" {
		return ""
	}

	var values map[string]*assignment.ResourceGroupValue
	if err := json.Unmarshal([]byte(jsonString.(string)), &values); err != nil {
		return fmt.Sprintf("unable to parse JSON: %+v", err)
	}

	b, _ := json.Marshal(values)
	return string(b)
}

func expandArmBlueprintAssignmentParameters(input string) map[string]assignment.ParameterValue {
	var result map[string]assignment.ParameterValue
	// the string has been validated by the schema, therefore the error is ignored here, since it will never happen.
	_ = json.Unmarshal([]byte(input), &result)
	return result
}

func expandArmBlueprintAssignmentResourceGroups(input string) map[string]assignment.ResourceGroupValue {
	var result map[string]assignment.ResourceGroupValue
	// the string has been validated by the schema, therefore the error is ignored here, since it will never happen.
	_ = json.Unmarshal([]byte(input), &result)
	return result
}

func flattenArmBlueprintAssignmentParameters(input map[string]assignment.ParameterValue) (string, error) {
	if len(input) == 0 {
		return "", nil
	}

	b, err := json.Marshal(input)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func flattenArmBlueprintAssignmentResourceGroups(input map[string]assignment.ResourceGroupValue) (string, error) {
	if len(input) == 0 {
		return "", nil
	}

	b, err := json.Marshal(input)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
