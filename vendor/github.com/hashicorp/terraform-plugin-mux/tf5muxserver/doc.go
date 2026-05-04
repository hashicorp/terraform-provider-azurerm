// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package tf5muxserver combines multiple provider servers that implement protocol version 5, into a single server.
//
// Supported protocol version 5 provider servers include any which implement
// the tfprotov5.ProviderServer (https://pkg.go.dev/github.com/hashicorp/terraform-plugin-go/tfprotov5#ProviderServer)
// interface, such as:
//
//   - https://pkg.go.dev/github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server
//   - https://pkg.go.dev/github.com/hashicorp/terraform-plugin-mux/tf6to5server
//   - https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema
//
// Refer to the NewMuxServer() function for creating a combined server.
package tf5muxserver
