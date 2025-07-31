// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// GetProviderSchemaResponse returns the *tfprotov5.GetProviderSchemaResponse
// equivalent of a *fwserver.GetProviderSchemaResponse.
func GetProviderSchemaResponse(ctx context.Context, fw *fwserver.GetProviderSchemaResponse) *tfprotov5.GetProviderSchemaResponse {
	if fw == nil {
		return nil
	}

	protov5 := &tfprotov5.GetProviderSchemaResponse{
		ActionSchemas:            make(map[string]*tfprotov5.ActionSchema, len(fw.ActionSchemas)),
		DataSourceSchemas:        make(map[string]*tfprotov5.Schema, len(fw.DataSourceSchemas)),
		Diagnostics:              Diagnostics(ctx, fw.Diagnostics),
		EphemeralResourceSchemas: make(map[string]*tfprotov5.Schema, len(fw.EphemeralResourceSchemas)),
		Functions:                make(map[string]*tfprotov5.Function, len(fw.FunctionDefinitions)),
		ListResourceSchemas:      make(map[string]*tfprotov5.Schema, len(fw.ListResourceSchemas)),
		ResourceSchemas:          make(map[string]*tfprotov5.Schema, len(fw.ResourceSchemas)),
		ServerCapabilities:       ServerCapabilities(ctx, fw.ServerCapabilities),
	}

	var err error

	protov5.Provider, err = Schema(ctx, fw.Provider)

	if err != nil {
		protov5.Diagnostics = append(protov5.Diagnostics, &tfprotov5.Diagnostic{
			Severity: tfprotov5.DiagnosticSeverityError,
			Summary:  "Error converting provider schema",
			Detail:   "The provider schema couldn't be converted into a usable type. This is always a problem with the provider. Please report the following to the provider developer:\n\n" + err.Error(),
		})

		return protov5
	}

	protov5.ProviderMeta, err = Schema(ctx, fw.ProviderMeta)

	if err != nil {
		protov5.Diagnostics = append(protov5.Diagnostics, &tfprotov5.Diagnostic{
			Severity: tfprotov5.DiagnosticSeverityError,
			Summary:  "Error converting provider_meta schema",
			Detail:   "The provider_meta schema couldn't be converted into a usable type. This is always a problem with the provider. Please report the following to the provider developer:\n\n" + err.Error(),
		})

		return protov5
	}

	for actionType, actionSchema := range fw.ActionSchemas {
		protov5.ActionSchemas[actionType], err = ActionSchema(ctx, actionSchema)

		if err != nil {
			protov5.Diagnostics = append(protov5.Diagnostics, &tfprotov5.Diagnostic{
				Severity: tfprotov5.DiagnosticSeverityError,
				Summary:  "Error converting action schema",
				Detail:   "The schema for the action \"" + actionType + "\" couldn't be converted into a usable type. This is always a problem with the provider. Please report the following to the provider developer:\n\n" + err.Error(),
			})

			return protov5
		}
	}

	for dataSourceType, dataSourceSchema := range fw.DataSourceSchemas {
		protov5.DataSourceSchemas[dataSourceType], err = Schema(ctx, dataSourceSchema)

		if err != nil {
			protov5.Diagnostics = append(protov5.Diagnostics, &tfprotov5.Diagnostic{
				Severity: tfprotov5.DiagnosticSeverityError,
				Summary:  "Error converting data source schema",
				Detail:   "The schema for the data source \"" + dataSourceType + "\" couldn't be converted into a usable type. This is always a problem with the provider. Please report the following to the provider developer:\n\n" + err.Error(),
			})

			return protov5
		}
	}

	for name, functionDefinition := range fw.FunctionDefinitions {
		protov5.Functions[name] = Function(ctx, functionDefinition)
	}

	for resourceType, resourceSchema := range fw.ResourceSchemas {
		protov5.ResourceSchemas[resourceType], err = Schema(ctx, resourceSchema)

		if err != nil {
			protov5.Diagnostics = append(protov5.Diagnostics, &tfprotov5.Diagnostic{
				Severity: tfprotov5.DiagnosticSeverityError,
				Summary:  "Error converting resource schema",
				Detail:   "The schema for the resource \"" + resourceType + "\" couldn't be converted into a usable type. This is always a problem with the provider. Please report the following to the provider developer:\n\n" + err.Error(),
			})

			return protov5
		}
	}

	for ephemeralResourceType, ephemeralResourceSchema := range fw.EphemeralResourceSchemas {
		protov5.EphemeralResourceSchemas[ephemeralResourceType], err = Schema(ctx, ephemeralResourceSchema)

		if err != nil {
			protov5.Diagnostics = append(protov5.Diagnostics, &tfprotov5.Diagnostic{
				Severity: tfprotov5.DiagnosticSeverityError,
				Summary:  "Error converting ephemeral resource schema",
				Detail:   "The schema for the ephemeral resource \"" + ephemeralResourceType + "\" couldn't be converted into a usable type. This is always a problem with the provider. Please report the following to the provider developer:\n\n" + err.Error(),
			})

			return protov5
		}
	}

	for listResourceType, listResourceSchema := range fw.ListResourceSchemas {
		protov5.ListResourceSchemas[listResourceType], err = Schema(ctx, listResourceSchema)

		if err != nil {
			protov5.Diagnostics = append(protov5.Diagnostics, &tfprotov5.Diagnostic{
				Severity: tfprotov5.DiagnosticSeverityError,
				Summary:  "Error converting list resource schema",
				Detail:   "The schema for the list resource \"" + listResourceType + "\" couldn't be converted into a usable type. This is always a problem with the provider. Please report the following to the provider developer:\n\n" + err.Error(),
			})
		}
	}

	return protov5
}
