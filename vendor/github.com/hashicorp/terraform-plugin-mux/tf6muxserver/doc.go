// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package tf6muxserver combines multiple provider servers that implement protocol version 6, into a single server.
//
// Supported protocol version 6 provider servers include any which implement
// the tfprotov6.ProviderServer (https://pkg.go.dev/github.com/hashicorp/terraform-plugin-go/tfprotov6#ProviderServer)
// interface, such as:
//
//   - https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework
//   - https://pkg.go.dev/github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server
//   - https://pkg.go.dev/github.com/hashicorp/terraform-plugin-mux/tf5to6server
//
// Refer to the NewMuxServer() function for creating a combined server.
package tf6muxserver
