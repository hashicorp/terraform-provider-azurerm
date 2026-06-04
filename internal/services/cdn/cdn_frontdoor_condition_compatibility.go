// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"fmt"

	"github.com/hashicorp/go-cty/cty"
)

func validateFrontDoorConditionBlocksRequireMatchValues(conditionBlock map[string]cty.Value, conditionNames []string) error {
	for _, conditionName := range conditionNames {
		conditionValue := conditionBlock[conditionName]
		if conditionValue.IsNull() || conditionValue.LengthInt() == 0 {
			continue
		}

		matchValues := conditionValue.AsValueSlice()[0].AsValueMap()["match_values"]
		if matchValues.IsNull() || matchValues.LengthInt() == 0 {
			return fmt.Errorf("the `%s` block requires `match_values`", conditionName)
		}
	}

	return nil
}
