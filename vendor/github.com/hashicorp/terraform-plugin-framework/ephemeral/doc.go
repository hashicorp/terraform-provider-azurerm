// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package ephemeral contains all interfaces, request types, and response
// types for an ephemeral resource implementation.
//
// In Terraform, an ephemeral resource is a concept which enables provider
// developers to offer practitioners ephemeral values, which will not be stored
// in any artifact produced by Terraform (plan/state). Ephemeral resources can
// optionally implement renewal logic via the (EphemeralResource).Renew method
// and cleanup logic via the (EphemeralResource).Close method.
//
// Ephemeral resources are not saved into the Terraform plan or state and can
// only be referenced in other ephemeral values, such as provider configuration
// attributes. Ephemeral resources are defined by a type/name, such as "examplecloud_thing",
// a schema representing the structure and data types of configuration, and lifecycle logic.
//
// The main starting point for implementations in this package is the
// EphemeralResource type which represents an instance of an ephemeral resource
// that has its own configuration and lifecycle logic. The [ephemeral.EphemeralResource]
// implementations are referenced by the [provider.ProviderWithEphemeralResources] type
// EphemeralResources method, which enables the ephemeral resource practitioner usage.
package ephemeral
