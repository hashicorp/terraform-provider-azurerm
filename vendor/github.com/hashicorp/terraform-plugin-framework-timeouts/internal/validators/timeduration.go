// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validators

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = timeDurationValidator{}

// timeDurationValidator validates that a string Attribute's value is parseable as time.Duration.
type timeDurationValidator struct {
}

// Description describes the validation in plain text formatting.
func (validator timeDurationValidator) Description(_ context.Context) string {
	return `must be a string containing a sequence of decimal numbers, each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".`
}

// MarkdownDescription describes the validation in Markdown formatting.
func (validator timeDurationValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

// ValidateString performs the validation.
func (validator timeDurationValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	s := req.ConfigValue

	if s.IsUnknown() || s.IsNull() {
		return
	}

	if _, err := time.ParseDuration(s.ValueString()); err != nil {
		resp.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(
			req.Path,
			"Invalid Attribute Value Time Duration",
			fmt.Sprintf("%q %s", s.ValueString(), validator.Description(ctx))),
		)
		return
	}
}

// TimeDuration returns an AttributeValidator which ensures that any configured
// attribute value:
//
//   - Is parseable as time duration.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func TimeDuration() validator.String {
	return timeDurationValidator{}
}
