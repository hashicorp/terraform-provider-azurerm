// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package provider contains all interfaces, request types, and response
// types for a provider implementation.
//
// In Terraform, a provider is a concept which enables provider developers
// to offer practitioners data sources and managed resources. Those concepts
// are described in more detail in their respective datasource and resource
// packages.
//
// Providers generally store any infrastructure clients or shared data that is
// applicable across data sources and managed resources. Providers are
// generally configured early in Terraform operations, such as plan and apply,
// before data source and managed resource logic is called. However, this early
// provider configuration is not guaranteed in the case there are unknown
// Terraform configuration values, so additional logic checks may be required
// throughout an implementation to handle this case. Providers may contain a
// schema representing the structure and data types of Terraform-based
// configuration.
//
// The main starting point for implementations in this package is the
// Provider type which represents an instance of a provider that has
// its own configuration.
package provider
