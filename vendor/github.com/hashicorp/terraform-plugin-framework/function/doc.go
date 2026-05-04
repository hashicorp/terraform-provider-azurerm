// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package function contains all interfaces, request types, and response
// types for a Terraform Provider function implementation.
//
// In Terraform, a function is a concept which enables provider developers
// to offer practitioners a pure function call in their configuration. Functions
// are defined by a function name, such as "parse_xyz", a definition
// representing the ordered list of parameters with associated data types and
// a result data type, and the function logic.
//
// The main starting point for implementations in this package is the
// [Function] type which represents an instance of a function that has its own
// argument data when called. The [Function] implementations are referenced by a
// [provider.Provider] type Functions method, which enables the function for
// practitioner and testing usage.
//
// Practitioner feedback is provided by the [FuncError] type, rather than
// the [diag.Diagnostic] type.
package function
