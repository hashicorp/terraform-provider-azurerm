// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package defaults contains schema default value interfaces and
// request/response implementations. These default value interfaces
// are used by resource/schema and internally in the framework.
// Refer to the typed default packages, such as stringdefault,
// for framework-defined default values that can be used in
// provider-defined schemas.
//
// Each attr.Type has a corresponding {TYPE} interface which
// implements concretely typed Default{TYPE} methods, such as
// StringDefault and DefaultString.
//
// The framework has to choose between default value developers handling a
// concrete framework value type, such as types.Bool, or the framework
// interface for custom value basetypes, such as basetypes.BoolValuable.
//
// In the framework type model, the developer can immediately use the value.
// If the value was associated with a custom type and using the custom value
// type is desired, the developer must use the type's ValueFrom{TYPE} method.
//
// In the custom type model, the developer must always convert to a concrete
// type before using the value unless checking for null or unknown. Since any
// custom type may be passed due to the schema, it is possible, if not likely,
// that unknown concrete types will be passed to the default.
//
// The framework chooses to pass the framework value type. This prevents the
// potential for unexpected runtime panics and simplifies development for
// easier use cases where the framework type is sufficient. More advanced
// developers can choose to call the type's ValueFrom{TYPE} method to get the
// desired custom type in a plan modifier.
//
// Defaults that are not type dependent need to implement all interfaces,
// but can use shared logic to reduce implementation code.
package defaults
