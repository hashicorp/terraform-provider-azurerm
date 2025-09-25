// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// GetProviderSchemaResponse returns the *tfprotov6.GetProviderSchemaResponse
// equivalent of a *fwserver.GetProviderSchemaResponse.
func GetProviderSchemaResponse(ctx context.Context, fw *fwserver.GetProviderSchemaResponse) *tfprotov6.GetProviderSchemaResponse {
	if fw == nil {
		return nil
	}

	protov6 := &tfprotov6.GetProviderSchemaResponse{
		ActionSchemas:            make(map[string]*tfprotov6.ActionSchema, len(fw.ActionSchemas)),
		DataSourceSchemas:        make(map[string]*tfprotov6.Schema, len(fw.DataSourceSchemas)),
		Diagnostics:              Diagnostics(ctx, fw.Diagnostics),
		EphemeralResourceSchemas: make(map[string]*tfprotov6.Schema, len(fw.EphemeralResourceSchemas)),
		Functions:                make(map[string]*tfprotov6.Function, len(fw.FunctionDefinitions)),
		ListResourceSchemas:      make(map[string]*tfprotov6.Schema, len(fw.ListResourceSchemas)),
		ResourceSchemas:          make(map[string]*tfprotov6.Schema, len(fw.ResourceSchemas)),
		ServerCapabilities:       ServerCapabilities(ctx, fw.ServerCapabilities),
	}

	var err error

	protov6.Provider, err = Schema(ctx, fw.Provider)

	if err != nil {
		protov6.Diagnostics = append(protov6.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error converting provider schema",
			Detail:   "The provider schema couldn't be converted into a usable type. This is always a problem with the provider. Please report the following to the provider developer:\n\n" + err.Error(),
		})

		return protov6
	}

	protov6.ProviderMeta, err = Schema(ctx, fw.ProviderMeta)

	if err != nil {
		protov6.Diagnostics = append(protov6.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error converting provider_meta schema",
			Detail:   "The provider_meta schema couldn't be converted into a usable type. This is always a problem with the provider. Please report the following to the provider developer:\n\n" + err.Error(),
		})

		return protov6
	}

	for actionType, actionSchema := range fw.ActionSchemas {
		protov6.ActionSchemas[actionType], err = ActionSchema(ctx, actionSchema)

		if err != nil {
			protov6.Diagnostics = append(protov6.Diagnostics, &tfprotov6.Diagnostic{
				Severity: tfprotov6.DiagnosticSeverityError,
				Summary:  "Error converting action schema",
				Detail:   "The schema for the action \"" + actionType + "\" couldn't be converted into a usable type. This is always a problem with the provider. Please report the following to the provider developer:\n\n" + err.Error(),
			})

			return protov6
		}
	}

	for dataSourceType, dataSourceSchema := range fw.DataSourceSchemas {
		protov6.DataSourceSchemas[dataSourceType], err = Schema(ctx, dataSourceSchema)

		if err != nil {
			protov6.Diagnostics = append(protov6.Diagnostics, &tfprotov6.Diagnostic{
				Severity: tfprotov6.DiagnosticSeverityError,
				Summary:  "Error converting data source schema",
				Detail:   "The schema for the data source \"" + dataSourceType + "\" couldn't be converted into a usable type. This is always a problem with the provider. Please report the following to the provider developer:\n\n" + err.Error(),
			})

			return protov6
		}
	}

	for name, functionDefinition := range fw.FunctionDefinitions {
		protov6.Functions[name] = Function(ctx, functionDefinition)
	}

	for resourceType, resourceSchema := range fw.ResourceSchemas {
		protov6.ResourceSchemas[resourceType], err = Schema(ctx, resourceSchema)

		if err != nil {
			protov6.Diagnostics = append(protov6.Diagnostics, &tfprotov6.Diagnostic{
				Severity: tfprotov6.DiagnosticSeverityError,
				Summary:  "Error converting resource schema",
				Detail:   "The schema for the resource \"" + resourceType + "\" couldn't be converted into a usable type. This is always a problem with the provider. Please report the following to the provider developer:\n\n" + err.Error(),
			})

			return protov6
		}
	}

	for ephemeralResourceType, ephemeralResourceSchema := range fw.EphemeralResourceSchemas {
		protov6.EphemeralResourceSchemas[ephemeralResourceType], err = Schema(ctx, ephemeralResourceSchema)

		if err != nil {
			protov6.Diagnostics = append(protov6.Diagnostics, &tfprotov6.Diagnostic{
				Severity: tfprotov6.DiagnosticSeverityError,
				Summary:  "Error converting ephemeral resource schema",
				Detail:   "The schema for the ephemeral resource \"" + ephemeralResourceType + "\" couldn't be converted into a usable type. This is always a problem with the provider. Please report the following to the provider developer:\n\n" + err.Error(),
			})

			return protov6
		}
	}

	for listResourceType, listResourceSchema := range fw.ListResourceSchemas {
		protov6.ListResourceSchemas[listResourceType], err = Schema(ctx, listResourceSchema)

		if err != nil {
			protov6.Diagnostics = append(protov6.Diagnostics, &tfprotov6.Diagnostic{
				Severity: tfprotov6.DiagnosticSeverityError,
				Summary:  "Error converting list resource schema",
				Detail:   "The schema for the list resource \"" + listResourceType + "\" couldn't be converted into a usable type. This is always a problem with the provider. Please report the following to the provider developer:\n\n" + err.Error(),
			})

			return protov6
		}
	}

	return protov6
}
