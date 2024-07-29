// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"errors"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/mitchellh/go-testing-interface"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
)

func runStateChecks(ctx context.Context, t testing.T, state *tfjson.State, stateChecks []statecheck.StateCheck) error {
	t.Helper()

	var result []error

	for _, stateCheck := range stateChecks {
		resp := statecheck.CheckStateResponse{}
		stateCheck.CheckState(ctx, statecheck.CheckStateRequest{State: state}, &resp)

		result = append(result, resp.Error)
	}

	return errors.Join(result...)
}
