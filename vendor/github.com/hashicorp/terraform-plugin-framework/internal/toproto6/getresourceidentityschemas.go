// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// GetResourceIdentitySchemasResponse returns the *tfprotov6.GetResourceIdentitySchemasResponse
// equivalent of a *fwserver.GetResourceIdentitySchemasResponse.
func GetResourceIdentitySchemasResponse(ctx context.Context, fw *fwserver.GetResourceIdentitySchemasResponse) *tfprotov6.GetResourceIdentitySchemasResponse {
	if fw == nil {
		return nil
	}

	protov6 := &tfprotov6.GetResourceIdentitySchemasResponse{
		Diagnostics:     Diagnostics(ctx, fw.Diagnostics),
		IdentitySchemas: make(map[string]*tfprotov6.ResourceIdentitySchema, len(fw.IdentitySchemas)),
	}

	var err error

	for resourceType, identitySchema := range fw.IdentitySchemas {
		protov6.IdentitySchemas[resourceType], err = IdentitySchema(ctx, identitySchema)

		if err != nil {
			protov6.Diagnostics = append(protov6.Diagnostics, &tfprotov6.Diagnostic{
				Severity: tfprotov6.DiagnosticSeverityError,
				Summary:  "Error converting resource identity schema",
				Detail:   "The identity schema for the resource \"" + resourceType + "\" couldn't be converted into a usable type. This is always a problem with the provider. Please report the following to the provider developer:\n\n" + err.Error(),
			})

			return protov6
		}
	}

	return protov6
}
