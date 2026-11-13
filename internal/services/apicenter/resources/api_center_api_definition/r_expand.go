// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_api_definition

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/apidefinitions"
)

func expandApiDefinitionPropertiesSpecification(input []ApiDefinitionPropertiesSpecification) *apidefinitions.ApiSpecImportRequest {
	if len(input) == 0 {
		return nil
	}
	v := input[0]

	if v.Url == "" && v.Content == "" {
		return nil
	}

	format := apidefinitions.ApiSpecImportSourceFormatInline
	content := v.Content
	if v.Url != "" {
		format = apidefinitions.ApiSpecImportSourceFormatLink
		content = v.Url
	}

	return &apidefinitions.ApiSpecImportRequest{
		Format: pointer.To(format),
		Specification: &apidefinitions.ApiSpecImportRequestSpecification{
			Name:    pointer.To(v.Name),
			Version: pointer.To(v.Version),
		},
		Value: pointer.To(content),
	}
}
