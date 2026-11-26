// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package consumption_budget_resource_group_test

import (
	"time"
)

func consumptionBudgetTestStartDate() time.Time {
	utcNow := time.Now().UTC()
	return time.Date(utcNow.Year(), utcNow.Month(), 1, 0, 0, 0, 0, utcNow.Location())
}
