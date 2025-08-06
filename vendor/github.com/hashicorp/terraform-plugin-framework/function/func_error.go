// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package function

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// NewFuncError returns a new function error with the
// given message.
func NewFuncError(text string) *FuncError {
	return &FuncError{
		Text: text,
	}
}

// NewArgumentFuncError returns a new function error with the
// given message and function argument.
func NewArgumentFuncError(functionArgument int64, text string) *FuncError {
	return &FuncError{
		Text:             text,
		FunctionArgument: &functionArgument,
	}
}

// FuncError is an error type specifically for function errors.
type FuncError struct {
	// Text is a practitioner-oriented description of the problem. This should
	// contain sufficient detail to provide both general and more specific information
	// regarding the issue. For example "Error executing function: foo can only contain
	// letters, numbers, and digits."
	Text string
	// FunctionArgument is a zero-based, int64 value that identifies the specific
	// function argument position that caused the error. Only errors that pertain
	// to a function argument will include this information.
	FunctionArgument *int64
}

// Equal returns true if the other function error is wholly equivalent.
func (fe *FuncError) Equal(other *FuncError) bool {
	if fe == nil && other == nil {
		return true
	}

	if fe == nil || other == nil {
		return false
	}

	if fe.Text != other.Text {
		return false
	}

	if fe.FunctionArgument == nil && other.FunctionArgument == nil {
		return true
	}

	if fe.FunctionArgument == nil || other.FunctionArgument == nil {
		return false
	}

	return *fe.FunctionArgument == *other.FunctionArgument
}

// Error returns the error text.
func (fe *FuncError) Error() string {
	if fe == nil {
		return ""
	}

	return fe.Text
}

// ConcatFuncErrors returns a new function error with the text from all supplied
// function errors concatenated together. If any of the function errors have a
// function argument, the first one encountered will be used.
func ConcatFuncErrors(funcErrs ...*FuncError) *FuncError {
	var text string
	var functionArgument *int64

	for _, f := range funcErrs {
		if f == nil {
			continue
		}

		if text != "" && f.Text != "" {
			text += "\n"
		}

		text += f.Text

		if functionArgument == nil {
			functionArgument = f.FunctionArgument
		}
	}

	if text != "" || functionArgument != nil {
		return &FuncError{
			Text:             text,
			FunctionArgument: functionArgument,
		}
	}

	return nil
}

// FuncErrorFromDiags iterates over the given diagnostics and returns a new function error
// with the summary and detail text from all error diagnostics concatenated together.
// Diagnostics with a severity of warning are logged but are not included in the returned
// function error.
func FuncErrorFromDiags(ctx context.Context, diags diag.Diagnostics) *FuncError {
	var funcErr *FuncError

	for _, d := range diags {
		switch d.Severity() {
		case diag.SeverityError:
			funcErr = ConcatFuncErrors(funcErr, NewFuncError(fmt.Sprintf("%s: %s", d.Summary(), d.Detail())))
		case diag.SeverityWarning:
			tflog.Warn(ctx, "warning: call function", map[string]interface{}{"summary": d.Summary(), "detail": d.Detail()})
		}
	}

	return funcErr
}
