// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// Implementation handler for a state move operation. After determining the
// source resource is supported, this should encapsulate all data transformation
// logic from a source resource to the current schema version of this
// [Resource]. The [Resource] is connected to these implementations by
// implementing the [ResourceWithMoveState] interface.
//
// This functionality is only available in Terraform 1.8 or later. It is invoked
// when a configuration contains a `moved` configuration block where the source
// resource type in the `from` argument differs from the target resource type in
// the `to` argument, causing Terraform to call this provider operation and the
// framework to route the request to this [Resource] as the target.
//
// Each implementation is responsible for determining whether the request should
// be handled or skipped. The implementation is considered skipped by the
// framework when the response contains no error diagnostics or state.
// Otherwise, any error diagnostics or state data, will cause the framework to
// consider the request completed and not call other [StateMover]
// implementations. The framework will return an implementation not found error
// if all implementations return responses without error diagnostics or state.
type StateMover struct {
	// SourceSchema is an optional schema for the intended source resource state
	// and schema version. While not required, setting this will populate
	// [MoveStateRequest.SourceState] when possible similar to other [Resource]
	// types.
	//
	// State conversion errors based on this schema are only logged at DEBUG
	// level as there may be multiple [StateMover] implementations on the same
	// target resource for differing source resources. The [StateMover]
	// implementation will still be called even with these errors, so it is
	// important that implementations verify the request via the
	// [MoveStateRequest.SourceTypeName] and other fields before attempting
	// to use [MoveStateRequest.SourceState].
	//
	// If not set, source state data is only available in
	// [MoveStateRequest.SourceRawState].
	SourceSchema *schema.Schema

	// StateMove defines the logic for determining whether the request source
	// resource information should match this implementation, and if so, the
	// data transformation of the source resource state to the current schema
	// version state of this [Resource].
	//
	// The [context.Context] parameter contains framework-defined loggers and
	// supports request cancellation.
	//
	// The [MoveStateRequest] parameter contains source resource information.
	// If [SourceSchema] was set, the [MoveStateRequest.SourceState] field will
	// be available. Otherwise, the [MoveStateRequest.SourceRawState] must be
	// used.
	//
	// The [MoveStateResponse] parameter can either remain unmodified to signal
	// to the framework that this implementation should be considered skipped or
	// must contain the transformed state data to signal a successful move. Any
	// returned error diagnostics will cause the framework to immediately
	// respond with those errors and without calling other [StateMover]
	// implementations.
	StateMover func(context.Context, MoveStateRequest, *MoveStateResponse)
}
