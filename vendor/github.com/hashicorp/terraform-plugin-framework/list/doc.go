// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package list contains all interfaces, request types, and response
// types for an list resource implementation.
//
// In Terraform, an list resource is a concept which enables provider
// developers to offer practitioners list values, which will not be stored
// in any artifact produced by Terraform (plan/state). List resources can
// optionally implement renewal logic via the (ListResource).Renew method
// and cleanup logic via the (ListResource).Close method.
//
// List resources are not saved into the Terraform plan or state and can
// only be referenced in other list values, such as provider configuration
// attributes. List resources are defined by a type/name, such as "examplecloud_thing",
// a schema representing the structure and data types of configuration, and lifecycle logic.
//
// The main starting point for implementations in this package is the
// ListResource type which represents an instance of an list resource
// that has its own configuration and lifecycle logic. The [list.ListResource]
// implementations are referenced by the [provider.ProviderWithListResources] type
// ListResources method, which enables the list resource practitioner usage.
package list
