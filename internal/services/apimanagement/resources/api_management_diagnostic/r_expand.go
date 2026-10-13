// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_diagnostic

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/diagnostic"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandApiManagementDiagnosticHTTPMessageDiagnostic(input []any) *diagnostic.HTTPMessageDiagnostic {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]any)

	result := &diagnostic.HTTPMessageDiagnostic{
		Body: &diagnostic.BodyDiagnosticSettings{},
	}

	if bodyBytes, ok := v["body_bytes"]; ok {
		result.Body.Bytes = pointer.To(int64(bodyBytes.(int)))
	}
	if headersSetRaw, ok := v["headers_to_log"]; ok {
		headersSet := headersSetRaw.(*pluginsdk.Set).List()
		headers := []string{}
		for _, header := range headersSet {
			headers = append(headers, header.(string))
		}
		result.Headers = pointer.To(headers)
	}

	result.DataMasking = expandApiManagementDataMasking(v["data_masking"].([]any))

	return result
}

func expandApiManagementDataMasking(input []any) *diagnostic.DataMasking {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	inputRaw := input[0].(map[string]any)
	return &diagnostic.DataMasking{
		QueryParams: expandApiManagementDataMaskingEntityList(inputRaw["query_params"].([]any)),
		Headers:     expandApiManagementDataMaskingEntityList(inputRaw["headers"].([]any)),
	}
}

func expandApiManagementDataMaskingEntityList(input []any) *[]diagnostic.DataMaskingEntity {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	result := make([]diagnostic.DataMaskingEntity, 0)
	for _, v := range input {
		entity := v.(map[string]any)
		result = append(result, diagnostic.DataMaskingEntity{
			Mode:  pointer.ToEnum[diagnostic.DataMaskingMode](entity["mode"].(string)),
			Value: pointer.To(entity["value"].(string)),
		})
	}
	return &result
}
