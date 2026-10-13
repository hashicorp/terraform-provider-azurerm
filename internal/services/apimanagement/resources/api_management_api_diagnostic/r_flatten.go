// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_api_diagnostic

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apidiagnostic"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/set"
)

func flattenApiManagementApiDiagnosticHTTPMessageDiagnostic(input *apidiagnostic.HTTPMessageDiagnostic) []any {
	result := make([]any, 0)

	if input == nil {
		return result
	}

	diagnostic := map[string]any{}

	if input.Body != nil && input.Body.Bytes != nil {
		diagnostic["body_bytes"] = pointer.From(input.Body.Bytes)
	}

	if input.Headers != nil {
		diagnostic["headers_to_log"] = set.FromStringSlice(*input.Headers)
	}

	diagnostic["data_masking"] = flattenApiManagementApiDiagnosticDataMasking(input.DataMasking)

	result = append(result, diagnostic)

	return result
}

func flattenApiManagementApiDiagnosticDataMasking(dataMasking *apidiagnostic.DataMasking) []any {
	if dataMasking == nil {
		return []any{}
	}

	var queryParams, headers []any
	if dataMasking.QueryParams != nil {
		queryParams = flattenApiManagementApiDiagnosticDataMaskingEntityList(dataMasking.QueryParams)
	}
	if dataMasking.Headers != nil {
		headers = flattenApiManagementApiDiagnosticDataMaskingEntityList(dataMasking.Headers)
	}

	return []any{
		map[string]any{
			"query_params": queryParams,
			"headers":      headers,
		},
	}
}

func flattenApiManagementApiDiagnosticDataMaskingEntityList(dataMaskingList *[]apidiagnostic.DataMaskingEntity) []any {
	if dataMaskingList == nil || len(*dataMaskingList) == 0 {
		return []any{}
	}

	result := []any{}

	for _, entity := range *dataMaskingList {
		result = append(result, map[string]any{
			"mode":  pointer.From(entity.Mode),
			"value": pointer.From(entity.Value),
		})
	}

	return result
}
