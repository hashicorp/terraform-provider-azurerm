// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_api_diagnostic

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apidiagnostic"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandApiManagementApiDiagnosticHTTPMessageDiagnostic(input []any) *apidiagnostic.HTTPMessageDiagnostic {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]any)

	result := &apidiagnostic.HTTPMessageDiagnostic{
		Body: &apidiagnostic.BodyDiagnosticSettings{},
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

	result.DataMasking = expandApiManagementApiDiagnosticDataMasking(v["data_masking"].([]any))

	return result
}

func expandApiManagementApiDiagnosticDataMasking(input []any) *apidiagnostic.DataMasking {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	inputRaw := input[0].(map[string]any)
	return &apidiagnostic.DataMasking{
		QueryParams: expandApiManagementApiDiagnosticDataMaskingEntityList(inputRaw["query_params"].([]any)),
		Headers:     expandApiManagementApiDiagnosticDataMaskingEntityList(inputRaw["headers"].([]any)),
	}
}

func expandApiManagementApiDiagnosticDataMaskingEntityList(input []any) *[]apidiagnostic.DataMaskingEntity {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	result := make([]apidiagnostic.DataMaskingEntity, 0)
	for _, v := range input {
		entity := v.(map[string]any)
		result = append(result, apidiagnostic.DataMaskingEntity{
			Mode:  pointer.ToEnum[apidiagnostic.DataMaskingMode](entity["mode"].(string)),
			Value: pointer.To(entity["value"].(string)),
		})
	}
	return &result
}
