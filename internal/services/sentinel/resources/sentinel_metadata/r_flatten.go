// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_metadata

import (
	"encoding/json"

	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/metadata"
)

func flattenMetadataSourceModel(input *metadata.MetadataSource) []MetadataSourceModel {
	if input == nil {
		return []MetadataSourceModel{}
	}
	output := MetadataSourceModel{
		Kind: string(input.Kind),
	}
	if input.Name != nil {
		output.Name = *input.Name
	}
	if input.SourceId != nil {
		output.Id = *input.SourceId
	}
	return []MetadataSourceModel{output}
}

func flattenMetadataAuthorModel(input *metadata.MetadataAuthor) []MetadataAuthorModel {
	if input == nil {
		return []MetadataAuthorModel{}
	}
	output := MetadataAuthorModel{}
	if input.Name != nil {
		output.Name = *input.Name
	}
	if input.Email != nil {
		output.Email = *input.Email
	}
	if input.Link != nil {
		output.Link = *input.Link
	}
	return []MetadataAuthorModel{output}
}

func flattenMetadataSupportModel(input *metadata.MetadataSupport) []MetadataSupportModel {
	if input == nil {
		return []MetadataSupportModel{}
	}
	output := MetadataSupportModel{
		Tier: string(input.Tier),
	}
	if input.Name != nil {
		output.Name = *input.Name
	}
	if input.Email != nil {
		output.Email = *input.Email
	}
	if input.Link != nil {
		output.Link = *input.Link
	}
	return []MetadataSupportModel{output}
}

func flattenMetadataCategoryModel(input *metadata.MetadataCategories) []MetadataCategoryModel {
	if input == nil {
		return []MetadataCategoryModel{}
	}
	output := MetadataCategoryModel{}
	if input.Domains != nil {
		output.Domains = *input.Domains
	}
	if input.Verticals != nil {
		output.Verticals = *input.Verticals
	}
	return []MetadataCategoryModel{output}
}

func flattenMetadataDependencies(input *metadata.MetadataDependencies) (string, error) {
	if input == nil {
		return "", nil
	}
	j, err := json.Marshal(input)
	if err != nil {
		return "", err
	}
	return string(j), nil
}
