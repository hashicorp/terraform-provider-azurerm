// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschemadata"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// ResourceIdentity returns the *tfsdk.ResourceIdentity for a *tfprotov6.ResourceIdentityData and fwschema.Schema.
func ResourceIdentity(ctx context.Context, in *tfprotov6.ResourceIdentityData, schema fwschema.Schema) (*tfsdk.ResourceIdentity, diag.Diagnostics) {
	if in == nil {
		return nil, nil
	}

	return IdentityData(ctx, in.IdentityData, schema)
}

// IdentityData returns the *tfsdk.ResourceIdentity for a *tfprotov6.DynamicValue and fwschema.Schema.
func IdentityData(ctx context.Context, proto6DynamicValue *tfprotov6.DynamicValue, schema fwschema.Schema) (*tfsdk.ResourceIdentity, diag.Diagnostics) {
	if proto6DynamicValue == nil {
		return nil, nil
	}

	var diags diag.Diagnostics

	// Panic prevention here to simplify the calling implementations.
	// This should not happen, but just in case.
	if schema == nil {
		diags.AddError(
			"Unable to Convert Resource Identity",
			"An unexpected error was encountered when converting the resource identity from the protocol type. "+
				"Identity data was sent in the protocol to a resource that doesn't support identity.\n\n"+
				"This is always a problem with Terraform or terraform-plugin-framework. Please report this to the provider developer.",
		)

		return nil, diags
	}

	data, dynamicValueDiags := DynamicValue(ctx, proto6DynamicValue, schema, fwschemadata.DataDescriptionResourceIdentity)

	diags.Append(dynamicValueDiags...)

	if diags.HasError() {
		return nil, diags
	}

	fw := &tfsdk.ResourceIdentity{
		Raw:    data.TerraformValue,
		Schema: schema,
	}

	return fw, diags
}
