// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/privatestate"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// MoveStateRequest represents a request for the provider to move a source
// resource state into the target resource state with any necessary data
// transformation logic. An instance of this request struct is supplied as an
// argument to each [StateMover.StateMover].
type MoveStateRequest struct {
	// SourcePrivate is the source resource private state data. If this source
	// data is important for the target resource, the implementation should set
	// the data via [MoveStateResponse.TargetPrivate] as it is intentionally not
	// copied automatically.
	SourcePrivate *privatestate.ProviderData

	// SourceProviderAddress is the address of the provider for the source
	// resource type. It is the full address in HOSTNAME/NAMESPACE/TYPE format.
	// For example, registry.terraform.io/hashicorp/random.
	//
	// Implementations should consider using this value to determine whether the
	// request should be handled by this particular implementation. It is
	// recommended to ignore the hostname unless necessary for disambiguation.
	SourceProviderAddress string

	// SourceRawState is the raw state of the source resource. This data is
	// always available, regardless whether the [StateMover.SourceSchema] field
	// was set. If SourceSchema is present, the [SourceState] field will be
	// populated and it is recommended to use that field instead.
	//
	// If this request matches the intended implementation, the implementation
	// logic must set [MoveStateResponse.State] as it is intentionally not
	// copied automatically.
	//
	// This is advanced functionality for providers wanting to skip the full
	// redeclaration of source schemas and instead use lower level handlers to
	// transform data. A typical implementation for working with this data will
	// call the Unmarshal() method.
	SourceRawState *tfprotov6.RawState

	// SourceSchemaVersion is the schema version of the source resource. It is
	// important for implementations to account for the schema version when
	// handling the source state data, since differing schema versions typically
	// have differing data structures and types.
	SourceSchemaVersion int64

	// SourceState is the source resource state if the [StateMover.SourceSchema]
	// was set. When available, this allows for easier data handling such as
	// calling Get() or GetAttribute().
	//
	// If this request matches the intended implementation, the implementation
	// logic must set [MoveStateResponse.TargetState] as it is intentionally not
	// copied automatically.
	//
	// State conversion errors based on [StateMover.SourceSchema] not matching
	// the source state are only intentionally logged at DEBUG level as there
	// may be multiple [StateMover] implementations on the same target resource
	// for differing source resources. The [StateMover] implementation will
	// still be called even with these errors, so it is important that
	// implementations verify the request via the SourceTypeName and other
	// fields before attempting to use this data.
	SourceState *tfsdk.State

	// SourceTypeName is the type name of the source resource. For example,
	// aws_vpc or random_string.
	//
	// Implementations should always use this value, in addition to potentially
	// other request fields, to determine whether the request should be handled
	// by this particular implementation.
	SourceTypeName string
}

// MoveStateResponse represents a response to a MoveStateRequest.
// An instance of this response struct is supplied as an argument to
// [StateMover] implementations. The provider should set response values only
// within the implementation that aligns with a supported request, or put
// differently, a response that contains no error diagnostics nor state is
// considered as skipped by the framework. Any fulfilling response, whether via
// error diagnostics or state data, will cause the framework to not call other
// implementations and immediately return that response. The framework will
// return an implementation not found error if all implementations return
// responses without error diagnostics and state.
type MoveStateResponse struct {
	// Diagnostics report errors or warnings related to moving the resource.
	// An unset or empty value indicates a successful operation or skipped
	// implementation if [TargetState] is also not set.
	Diagnostics diag.Diagnostics

	// TargetState is the resource state following the move operation. This
	// value is intentionally not pre-populated by the framework. The provider
	// must set state values to indicate a successful move operation.
	//
	// If this value is unset and there are no diagnostics, the framework will
	// consider this implementation to have not matched the request and move on
	// to the next implementation, if any. If no implementation returns a state
	// then the framework will return an implementation not found error.
	TargetState tfsdk.State

	// TargetPrivate is the resource private state data following the move
	// operation. This field is not pre-populated as it is up to implementations
	// whether the source private data is relevant for the target resource.
	TargetPrivate *privatestate.ProviderData
}
