// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package datasource contains all interfaces, request types, and response
// types for a data source implementation.
//
// In Terraform, a data source is a concept which enables provider developers
// to offer practitioners a read-only source of information, which is saved
// into the Terraform state and can be referenced by other parts of a
// configuration. Data sources are defined by a data source type/name, such as
// "examplecloud_thing", a schema representing the structure and data types of
// configuration and state, and read logic.
//
// The main starting point for implementations in this package is the
// DataSource type which represents an instance of a data source type that has
// its own configuration, read logic, and state. The DataSource implementations
// are referenced by a [provider.Provider] type DataSources method, which
// enables the data source for practitioner and testing usage.
package datasource
