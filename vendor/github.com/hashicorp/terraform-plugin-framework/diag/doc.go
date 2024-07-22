// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package diag implements diagnostic functionality, which is a practitioner
// feedback mechanism for providers. It is designed for display in Terraform
// user interfaces, rather than logging based feedback, which is generally
// saved to a file for later inspection and troubleshooting.
//
// Practitioner feedback for provider defined functions is provided by the
// [function.FuncError] type, rather than the [diag.Diagnostic] type.
package diag
