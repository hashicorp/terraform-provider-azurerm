// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package types contains the framework-defined data types and values, such as
// boolean, floating point, integer, list, map, object, set, and string.
//
// This package contains creation functions and type aliases for most provider
// use cases. The actual schema-ready type and value type implementations are
// under the basetypes package. Embed those basetypes implementations to create
// custom types.
package types
