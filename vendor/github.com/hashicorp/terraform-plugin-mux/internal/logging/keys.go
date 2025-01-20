// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logging

// Global logging keys attached to all requests.
//
// Practitioners or tooling reading logs may be depending on these keys, so be
// conscious of that when changing them.
const (
	// Go type of the provider selected by mux.
	KeyTfMuxProvider = "tf_mux_provider"

	// The RPC being run, such as "ApplyResourceChange"
	KeyTfRpc = "tf_rpc"
)
