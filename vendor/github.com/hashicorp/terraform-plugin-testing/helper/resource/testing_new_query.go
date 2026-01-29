// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/mitchellh/go-testing-interface"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource/query"
	"github.com/hashicorp/terraform-plugin-testing/internal/plugintest"
	"github.com/hashicorp/terraform-plugin-testing/internal/teststep"
)

func testStepNewQuery(ctx context.Context, t testing.T, wd *plugintest.WorkingDir, step TestStep, providers *providerFactories) error {
	t.Helper()

	queryConfigRequest := teststep.ConfigurationRequest{
		Raw: &step.Config,
	}
	err := wd.SetQuery(ctx, teststep.Configuration(queryConfigRequest), step.ConfigVariables)
	if err != nil {
		return fmt.Errorf("Error setting query config: %w", err)
	}

	err = runProviderCommand(ctx, t, wd, providers, func() error {
		return wd.Init(ctx)
	})
	if err != nil {
		t.Fatalf("Error getting init: %s", err)
	}

	var queryOut []tfjson.LogMsg
	err = runProviderCommand(ctx, t, wd, providers, func() error {
		var err error
		queryOut, err = wd.Query(ctx)
		return err
	})
	if err != nil {
		return err
	}

	return query.RunQueryChecks(ctx, t, queryOut, step.QueryResultChecks)
}
