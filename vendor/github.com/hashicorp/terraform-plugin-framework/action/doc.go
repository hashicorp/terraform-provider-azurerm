// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package action contains all interfaces, request types, and response
// types for an action implementation.
//
// In Terraform, an action is a concept which enables provider developers
// to offer practitioners ad-hoc side-effects to be used in their configuration.
//
// The main starting point for implementations in this package is the
// [Action] type which represents an instance of an action that has its
// own configuration, plan, and invoke logic. The [Action] implementations
// are referenced by the [provider.ProviderWithActions] type Actions method,
// which enables the action practitioner usage.
package action
