// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package function

// RunRequest represents a request for the Function to call its implementation
// logic. An instance of this request struct is supplied as an argument to the
// Function type Run method.
type RunRequest struct {
	// Arguments is the data sent from Terraform. Use the ArgumentsData type
	// GetArgument method to retrieve each positional argument.
	Arguments ArgumentsData
}

// RunResponse represents a response to a RunRequest. An instance of this
// response struct is supplied as an argument to the Function type Run method.
type RunResponse struct {
	// Error contains errors related to running the function.
	// A nil error indicates success, with no errors generated.
	// [ConcatFuncErrors] can be used to combine multiple errors into a single error.
	Error *FuncError

	// Result is the data to be returned to Terraform matching the function
	// result definition. This must be set or an error diagnostic is raised. Use
	// the ResultData type Set method to save the data.
	Result ResultData
}
