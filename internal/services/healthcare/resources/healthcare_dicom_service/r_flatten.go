// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package healthcare_dicom_service

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2024-03-31/dicomservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func flattenDicomAuthentication(input *dicomservices.DicomServiceAuthenticationConfiguration) []any {
	if input == nil {
		return []any{}
	}

	authBlock := make(map[string]any)
	if input.Authority != nil {
		authBlock["authority"] = *input.Authority
	}

	audience := make([]any, 0)
	if input.Audiences != nil {
		for _, data := range *input.Audiences {
			audience = append(audience, data)
		}
	}
	authBlock["audience"] = audience

	return []any{authBlock}
}

func flattenDicomServicePrivateEndpoint(input *[]dicomservices.PrivateEndpointConnection) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, endpoint := range *input {
		result := map[string]any{}
		if endpoint.Name != nil {
			result["name"] = *endpoint.Name
		}

		if endpoint.Id != nil {
			result["id"] = *endpoint.Id
		}
	}
	return results
}

func flattenStorageConfiguration(configuration *dicomservices.StorageConfiguration) (any, error) {
	if configuration == nil {
		return []any{}, nil
	}

	result := make(map[string]any)
	if configuration.FileSystemName != nil {
		result["file_system_name"] = pointer.From(configuration.FileSystemName)
	}

	if v := pointer.From(configuration.StorageResourceId); v != "" {
		id, err := commonids.ParseStorageAccountID(v)
		if err != nil {
			return nil, err
		}
		result["storage_account_id"] = id.ID()
	}

	return []any{result}, nil
}

func flattenDicomServiceCorsConfiguration(input *dicomservices.CorsConfiguration) []any {
	outputList := make([]any, 0)
	if input == nil {
		return outputList
	}

	output := make(map[string]any)
	output["allow_credentials"] = pointer.From(input.AllowCredentials)

	if input.Headers != nil {
		output["allowed_headers"] = pluginsdk.FlattenSlice(input.Headers)
	}

	output["max_age_in_seconds"] = pointer.From(input.MaxAge)

	if input.Methods != nil {
		output["allowed_methods"] = pluginsdk.FlattenSlice(input.Methods)
	}

	if input.Origins != nil {
		output["allowed_origins"] = pluginsdk.FlattenSlice(input.Origins)
	}

	return append(outputList, output)
}
