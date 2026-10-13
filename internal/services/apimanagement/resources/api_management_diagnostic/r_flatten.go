// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_diagnostic

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/diagnostic"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/set"
)

func flattenApiManagementDiagnosticHTTPMessageDiagnostic(input *diagnostic.HTTPMessageDiagnostic) []any {
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

	diagnostic["data_masking"] = flattenApiManagementDataMasking(input.DataMasking)

	result = append(result, diagnostic)

	return result
}

func flattenApiManagementDataMasking(dataMasking *diagnostic.DataMasking) []any {
	if dataMasking == nil {
		return []any{}
	}

	var queryParams, headers []any
	if dataMasking.QueryParams != nil {
		queryParams = flattenApiManagementDataMaskingEntityList(dataMasking.QueryParams)
	}
	if dataMasking.Headers != nil {
		headers = flattenApiManagementDataMaskingEntityList(dataMasking.Headers)
	}

	return []any{
		map[string]any{
			"query_params": queryParams,
			"headers":      headers,
		},
	}
}

func flattenApiManagementDataMaskingEntityList(dataMaskingList *[]diagnostic.DataMaskingEntity) []any {
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
