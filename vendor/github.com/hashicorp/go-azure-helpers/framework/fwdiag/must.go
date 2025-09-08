// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwdiag

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// Must - ensures a hard fail if diags contains errors from the supplied function x. This should
// protect developers from shipping panics from incompatible type conversions /
func Must[T any](x T, diags diag.Diagnostics) T {
	return ErrMust(x, DiagAsError(diags))
}

func ErrMust[T any](x T, err error) T {
	if err != nil {
		panic(err)
	}
	return x
}

func DiagAsError(diags diag.Diagnostics) error {
	errs := make([]error, 0)

	for _, err := range diags.Errors() {
		errStr := err.Summary()
		if err.Detail() != "" {
			errStr += ": " + err.Detail()
		}
		errs = append(errs, errors.New(errStr))
	}

	return errors.Join(errs...)
}
