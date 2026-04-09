// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwschema

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

// NumericPrefixRegex is a regular expression which matches whether a string
// begins with a numeric (0-9).
var NumericPrefixRegex = regexp.MustCompile(`^[0-9]`)

// ReservedProviderAttributeNames contains the list of root attribute names
// which should not be included in provider-defined provider schemas since
// they require practitioners to implement special syntax in their
// configurations to be usable by the provider.
var ReservedProviderAttributeNames = []string{
	// Reference: https://developer.hashicorp.com/terraform/language/providers/configuration#alias-multiple-provider-configurations
	"alias",
	// Reference: https://developer.hashicorp.com/terraform/language/providers/configuration#version-deprecated
	"version",
}

// ReservedResourceAttributeNames contains the list of root attribute names
// which should not be included in provider-defined managed resource and
// data source schemas since they require practitioners to implement special
// syntax in their configurations to be usable by the provider resource.
var ReservedResourceAttributeNames = []string{
	// Reference: https://developer.hashicorp.com/terraform/language/resources/provisioners/connection
	"connection",
	// Reference: https://developer.hashicorp.com/terraform/language/meta-arguments/count
	"count",
	// Reference: https://developer.hashicorp.com/terraform/language/meta-arguments/depends_on
	"depends_on",
	// Reference: https://developer.hashicorp.com/terraform/language/meta-arguments/for_each
	"for_each",
	// Reference: https://developer.hashicorp.com/terraform/language/meta-arguments/lifecycle
	"lifecycle",
	// Reference: https://developer.hashicorp.com/terraform/language/meta-arguments/resource-provider
	"provider",
	// Reference: https://developer.hashicorp.com/terraform/language/resources/provisioners/syntax
	"provisioner",
}

// ValidAttributeNameRegex contains the regular expression to validate
// attribute names, which are considered [identifiers] in the Terraform
// configuration language.
//
// Hyphen characters (-) are technically valid in identifiers, however they are
// explicitly not validated due to the provider ecosystem conventionally never
// including them in attribute names. Introducing them could cause practitioner
// confusion.
//
// [identifiers]: https://developer.hashicorp.com/terraform/language/syntax/configuration#identifiers
var ValidAttributeNameRegex = regexp.MustCompile("^[a-z_][a-z0-9_]*$")

// IsReservedProviderAttributeName returns an error diagnostic if the given
// attribute path represents a root attribute name in
// ReservedProviderAttributeNames. Other paths are automatically skipped
// without error.
func IsReservedProviderAttributeName(name string, attributePath path.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Check that given path is root attribute name. This simplifies calling
	// logic to not worry about conditionalizing this check.
	if len(attributePath.Steps()) != 1 {
		return diags
	}

	for _, reservedName := range ReservedProviderAttributeNames {
		if name == reservedName {
			// The diagnostic path is intentionally omitted as it is invalid
			// in this context. Diagnostic paths are intended to be mapped to
			// actual data, while this path information must be synthesized.
			diags.AddError(
				"Reserved Root Attribute/Block Name",
				"When validating the provider schema, an implementation issue was found. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("%q is a reserved root attribute/block name. ", name)+
					"This is to prevent practitioners from needing special Terraform configuration syntax.",
			)

			break
		}
	}

	return diags
}

// IsReservedResourceAttributeName returns an error diagnostic if the given
// attribute path represents a root attribute name in
// ReservedResourceAttributeNames. Other paths are automatically skipped
// without error.
func IsReservedResourceAttributeName(name string, attributePath path.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Check that given path is root attribute name. This simplifies calling
	// logic to not worry about conditionalizing this check.
	if len(attributePath.Steps()) != 1 {
		return diags
	}

	for _, reservedName := range ReservedResourceAttributeNames {
		if name == reservedName {
			// The diagnostic path is intentionally omitted as it is invalid
			// in this context. Diagnostic paths are intended to be mapped to
			// actual data, while this path information must be synthesized.
			diags.AddError(
				"Reserved Root Attribute/Block Name",
				"When validating the resource or data source schema, an implementation issue was found. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("%q is a reserved root attribute/block name. ", name)+
					"This is to prevent practitioners from needing special Terraform configuration syntax.",
			)

			break
		}
	}

	return diags
}

// IsValidAttributeName returns an error diagnostic if the given
// attribute path has an invalid attribute name according to
// ValidAttributeNameRegex. Non-AttributeName paths are automatically skipped
// without error.
func IsValidAttributeName(name string, attributePath path.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	if ValidAttributeNameRegex.MatchString(name) {
		return diags
	}

	var message strings.Builder

	message.WriteString("Names must ")

	if NumericPrefixRegex.MatchString(name) {
		message.WriteString("begin with a lowercase alphabet character (a-z) or underscore (_) and must ")
	}

	message.WriteString("only contain lowercase alphanumeric characters (a-z, 0-9) and underscores (_).")

	// The diagnostic path is intentionally omitted as it is invalid in this
	// context. Diagnostic paths are intended to be mapped to actual data,
	// while this path information must be synthesized.
	diags.AddError(
		"Invalid Attribute/Block Name",
		"When validating the schema, an implementation issue was found. "+
			"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
			fmt.Sprintf("%q at schema path %q is an invalid attribute/block name. ", name, attributePath)+
			message.String(),
	)

	return diags
}
