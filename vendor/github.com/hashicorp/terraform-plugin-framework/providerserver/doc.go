// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package providerserver implements functionality for serving a provider,
// such as directly starting a server in a production binary and conversion
// functions for testing.
//
// For production usage, call the Serve function from binary startup, such as
// from the provider codebase main package. If multiplexing the provider server
// via terraform-plugin-mux functionality, use the NewProtocol* functions and
// call the Serve function from that Go module. For testing usage, call the
// NewProtocol* functions.
//
// All functionality in this package requires the provider.Provider type, which
// contains the provider implementation including all managed resources and
// data sources.
package providerserver
