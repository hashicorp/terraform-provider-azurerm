// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// GetResourceIdentitySchemasResponse returns the *tfprotov5.GetResourceIdentitySchemasResponse
// equivalent of a *fwserver.GetResourceIdentitySchemasResponse.
func GetResourceIdentitySchemasResponse(ctx context.Context, fw *fwserver.GetResourceIdentitySchemasResponse) *tfprotov5.GetResourceIdentitySchemasResponse {
	if fw == nil {
		return nil
	}

	protov5 := &tfprotov5.GetResourceIdentitySchemasResponse{
		Diagnostics:     Diagnostics(ctx, fw.Diagnostics),
		IdentitySchemas: make(map[string]*tfprotov5.ResourceIdentitySchema, len(fw.IdentitySchemas)),
	}

	var err error

	for resourceType, identitySchema := range fw.IdentitySchemas {
		protov5.IdentitySchemas[resourceType], err = IdentitySchema(ctx, identitySchema)

		if err != nil {
			protov5.Diagnostics = append(protov5.Diagnostics, &tfprotov5.Diagnostic{
				Severity: tfprotov5.DiagnosticSeverityError,
				Summary:  "Error converting resource identity schema",
				Detail:   "The identity schema for the resource \"" + resourceType + "\" couldn't be converted into a usable type. This is always a problem with the provider. Please report the following to the provider developer:\n\n" + err.Error(),
			})

			return protov5
		}
	}

	return protov5
}
