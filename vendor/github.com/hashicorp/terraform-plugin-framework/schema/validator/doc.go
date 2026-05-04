// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package validator contains common schema validator interfaces and
// implementations. These validators are used by concept specific packages
// such as datasource/schema, provider/schema, and resource/schema.
//
// Each attr.Type has a corresponding {TYPE}Validator interface which
// implements concretely typed Validate{TYPE} methods, such as
// StringValidator and ValidateString. Custom attr.Type can also consider
// implementing native type validation via the attr/xattr.TypeWithValidate
// interface instead of schema validators.
//
// The framework has to choose between validator developers handling a concrete
// framework value type, such as types.Bool, or the framework interface for
// custom value basetypes. such as basetypes.BoolValuable.
//
// In the framework type model, the developer can immediately use the value.
// If the value was associated with a custom type and using the custom value
// type is desired, the developer must use the type's ValueFrom{TYPE} method.
//
// In the custom type model, the developer must always convert to a concreate
// type before using the value unless checking for null or unknown. Since any
// custom type may be passed due to the schema, it is possible, if not likely,
// that unknown concrete types will be passed to the validator.
//
// The framework chooses to pass the framework value type. This prevents the
// potential for unexpected runtime panics and simplifies development for
// easier use cases where the framework type is sufficient. More advanced
// developers can choose to implement native type validation for custom
// types or call the type's ValueFrom{TYPE} method to get the desired
// desired custom type in a validator.
//
// Validators that are not type dependent need to implement all interfaces,
// but can use shared logic to reduce implementation code.
package validator
