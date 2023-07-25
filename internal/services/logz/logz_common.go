// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logz

import "fmt"

const PlanId100gb14days = "100gb14days"
const PlanDetails100gb14days = "100gb14days@TIDgmz7xq9ge3py"

func getPlanDetails(plan string) (string, error) {
	if plan == PlanId100gb14days {
		return PlanDetails100gb14days, nil
	}

	return "", fmt.Errorf("cannot find plan details for id: %s", plan)
}

func getPlanId(planDetails string) (string, error) {
	if planDetails == PlanDetails100gb14days {
		return PlanId100gb14days, nil
	}

	return "", fmt.Errorf("cannot find plan id for details: %s", planDetails)
}
