// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfprotov6

// ActionSchema is how Terraform defines the shape of action data.
type ActionSchema struct {
	// Schema is the definition for the action data itself, which will be specified in an action block in the user's configuration.
	Schema *Schema
}
