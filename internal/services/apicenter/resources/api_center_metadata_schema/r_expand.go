// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_metadata_schema

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/metadataschemas"
)

func expandMetadataAssignmentList(input []MetadataAssignment) *[]metadataschemas.MetadataAssignment {
	if len(input) == 0 {
		return nil
	}
	result := make([]metadataschemas.MetadataAssignment, 0, len(input))
	for _, v := range input {
		result = append(result, metadataschemas.MetadataAssignment{
			Deprecated: pointer.To(v.Deprecated),
			Entity:     pointer.ToEnum[metadataschemas.MetadataAssignmentEntity](v.Entity),
			Required:   pointer.To(v.Required),
		})
	}
	return &result
}
