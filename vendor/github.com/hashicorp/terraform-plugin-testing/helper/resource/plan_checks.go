// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/internal/errorshim"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/mitchellh/go-testing-interface"
)

func runPlanChecks(ctx context.Context, t testing.T, plan *tfjson.Plan, planChecks []plancheck.PlanCheck) error {
	t.Helper()

	var result error

	for _, planCheck := range planChecks {
		resp := plancheck.CheckPlanResponse{}
		planCheck.CheckPlan(ctx, plancheck.CheckPlanRequest{Plan: plan}, &resp)

		if resp.Error != nil {
			// TODO: Once Go 1.20 is the minimum supported version for this module, replace with `errors.Join` function
			// - https://github.com/hashicorp/terraform-plugin-testing/issues/99
			result = errorshim.Join(result, resp.Error)
		}
	}

	return result
}
