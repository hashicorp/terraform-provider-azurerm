// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package schemavalidator provides validators to express relationships between
// multiple attributes within the schema of a resource, data source, or provider.
// For example, checking that an attribute is present when another is present, or vice-versa.
//
// These validators are implemented on a starting attribute, where
// relationships can be expressed as absolute paths to others or relative to
// the starting attribute. For multiple attribute validators that are defined
// outside the schema, which may be easier to implement in provider code
// generation situations or suit provider code preferences differently, refer
// to the datasourcevalidator, providervalidator, or resourcevalidator package.
package schemavalidator
