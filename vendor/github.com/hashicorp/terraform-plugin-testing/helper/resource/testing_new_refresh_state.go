// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/mitchellh/go-testing-interface"

	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-plugin-testing/internal/logging"
	"github.com/hashicorp/terraform-plugin-testing/internal/plugintest"
)

func testStepNewRefreshState(ctx context.Context, t testing.T, wd *plugintest.WorkingDir, step TestStep, providers *providerFactories) error {
	t.Helper()

	var err error
	// Explicitly ensure prior state exists before refresh.
	err = runProviderCommand(ctx, t, func() error {
		_, err = getState(ctx, t, wd)
		if err != nil {
			return err
		}
		return nil
	}, wd, providers)
	if err != nil {
		t.Fatalf("Error getting state: %s", err)
	}

	err = runProviderCommand(ctx, t, func() error {
		return wd.Refresh(ctx)
	}, wd, providers)
	if err != nil {
		return err
	}

	var refreshState *terraform.State
	err = runProviderCommand(ctx, t, func() error {
		refreshState, err = getState(ctx, t, wd)
		if err != nil {
			return err
		}
		return nil
	}, wd, providers)
	if err != nil {
		t.Fatalf("Error getting state: %s", err)
	}

	// Go through the refreshed state and verify
	if step.Check != nil {
		logging.HelperResourceDebug(ctx, "Calling TestStep Check for RefreshState")

		if err := step.Check(refreshState); err != nil {
			t.Fatal(err)
		}

		logging.HelperResourceDebug(ctx, "Called TestStep Check for RefreshState")
	}

	// do a plan
	err = runProviderCommand(ctx, t, func() error {
		return wd.CreatePlan(ctx)
	}, wd, providers)
	if err != nil {
		return fmt.Errorf("Error running post-refresh plan: %w", err)
	}

	var plan *tfjson.Plan
	err = runProviderCommand(ctx, t, func() error {
		var err error
		plan, err = wd.SavedPlan(ctx)
		return err
	}, wd, providers)
	if err != nil {
		return fmt.Errorf("Error retrieving post-refresh plan: %w", err)
	}

	// Run post-refresh plan checks
	if len(step.RefreshPlanChecks.PostRefresh) > 0 {
		err = runPlanChecks(ctx, t, plan, step.RefreshPlanChecks.PostRefresh)
		if err != nil {
			return fmt.Errorf("Post-refresh plan check(s) failed:\n%w", err)
		}
	}

	if !planIsEmpty(plan, wd.GetHelper().TerraformVersion()) && !step.ExpectNonEmptyPlan {
		var stdout string
		err = runProviderCommand(ctx, t, func() error {
			var err error
			stdout, err = wd.SavedPlanRawStdout(ctx)
			return err
		}, wd, providers)
		if err != nil {
			return fmt.Errorf("Error retrieving formatted plan output: %w", err)
		}
		return fmt.Errorf("After refreshing state during this test step, a followup plan was not empty.\nstdout:\n\n%s", stdout)
	}

	return nil
}
