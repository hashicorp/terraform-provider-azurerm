// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plancheck

// ResourceActionType is a string enum type that routes to a specific terraform-json.Actions function for asserting resource changes.
//   - https://pkg.go.dev/github.com/hashicorp/terraform-json#Actions
//
// More information about expected resource behavior can be found at: https://developer.hashicorp.com/terraform/language/resources/behavior
type ResourceActionType string

const (
	// ResourceActionNoop occurs when a resource is not planned to change (no-op).
	//   - Routes to: https://pkg.go.dev/github.com/hashicorp/terraform-json#Actions.NoOp
	ResourceActionNoop ResourceActionType = "NoOp"

	// ResourceActionCreate occurs when a resource is planned to be created.
	//   - Routes to: https://pkg.go.dev/github.com/hashicorp/terraform-json#Actions.Create
	ResourceActionCreate ResourceActionType = "Create"

	// ResourceActionRead occurs when a data source is planned to be read during the apply stage (data sources are read during plan stage when possible).
	// See the data source documentation for more information on this behavior: https://developer.hashicorp.com/terraform/language/data-sources#data-resource-behavior
	//   - Routes to: https://pkg.go.dev/github.com/hashicorp/terraform-json#Actions.Read
	ResourceActionRead ResourceActionType = "Read"

	// ResourceActionUpdate occurs when a resource is planned to be updated in-place.
	//   - Routes to: https://pkg.go.dev/github.com/hashicorp/terraform-json#Actions.Update
	ResourceActionUpdate ResourceActionType = "Update"

	// ResourceActionDestroy occurs when a resource is planned to be deleted.
	//   - Routes to: https://pkg.go.dev/github.com/hashicorp/terraform-json#Actions.Delete
	ResourceActionDestroy ResourceActionType = "Destroy"

	// ResourceActionDestroyBeforeCreate occurs when a resource is planned to be deleted and then re-created. This is the default
	// behavior when terraform must change a resource argument that cannot be updated in-place due to remote API limitations.
	//   - Routes to: https://pkg.go.dev/github.com/hashicorp/terraform-json#Actions.DestroyBeforeCreate
	ResourceActionDestroyBeforeCreate ResourceActionType = "DestroyBeforeCreate"

	// ResourceActionCreateBeforeDestroy occurs when a resource is planned to be created and then deleted. This is opt-in behavior that
	// is enabled with the [create_before_destroy] meta-argument.
	//   - Routes to: https://pkg.go.dev/github.com/hashicorp/terraform-json#Actions.CreateBeforeDestroy
	//
	// [create_before_destroy]: https://developer.hashicorp.com/terraform/language/meta-arguments/lifecycle#create_before_destroy
	ResourceActionCreateBeforeDestroy ResourceActionType = "CreateBeforeDestroy"

	// ResourceActionReplace can be used to verify a resource is planned to be deleted and re-created (where the order of delete and create actions are not important).
	// This action matches both ResourceActionDestroyBeforeCreate and ResourceActionCreateBeforeDestroy.
	//   - Routes to: https://pkg.go.dev/github.com/hashicorp/terraform-json#Actions.Replace
	ResourceActionReplace ResourceActionType = "Replace"
)
