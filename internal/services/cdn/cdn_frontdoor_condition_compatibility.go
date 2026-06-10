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
		if conditionValue.IsNull() || !conditionValue.IsKnown() || conditionValue.LengthInt() == 0 {
			continue
		}

		conditionEntries := conditionValue.AsValueSlice()
		if len(conditionEntries) == 0 || conditionEntries[0].IsNull() || !conditionEntries[0].IsKnown() {
			continue
		}

		matchValues := conditionEntries[0].AsValueMap()["match_values"]
		if matchValues.IsNull() || !matchValues.IsKnown() {
			continue
		}

		if matchValues.LengthInt() == 0 {
			return fmt.Errorf("the `%s` block requires `match_values`", conditionName)
		}
	}

	return nil
}
