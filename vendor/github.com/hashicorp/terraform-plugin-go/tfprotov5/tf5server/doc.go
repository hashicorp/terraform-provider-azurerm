// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package tf5server implements a server implementation to run
// tfprotov5.ProviderServers as gRPC servers.
//
// Providers will likely be calling tf5server.Serve from their main function to
// start the server so Terraform can connect to it.
package tf5server
