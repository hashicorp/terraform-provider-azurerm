// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights_workbook_template

import (
	"encoding/json"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-11-20/workbooktemplatesapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *workbooktemplatesapis.WorkbookTemplateId, model *workbooktemplatesapis.WorkbookTemplate) error {
	state := ApplicationInsightsWorkbookTemplateModel{
		Name:              id.WorkbookTemplateName,
		ResourceGroupName: id.ResourceGroupName,
		Location:          location.Normalize(model.Location),
	}

	if properties := model.Properties; properties != nil {
		if properties.Author != nil {
			state.Author = *properties.Author
		}

		state.Galleries = flattenWorkbookTemplateGalleryModel(&properties.Galleries)

		if properties.Priority != nil {
			state.Priority = *properties.Priority
		}

		if properties.TemplateData != nil {
			templateDataValue, err := json.Marshal(properties.TemplateData)
			if err != nil {
				return err
			}

			state.TemplateData = string(templateDataValue)
		}

		if properties.Localized != nil {
			localizedValue, err := json.Marshal(properties.Localized)
			if err != nil {
				return err
			}

			state.Localized = string(localizedValue)
		}
	}

	if model.Tags != nil {
		state.Tags = *model.Tags
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func flattenWorkbookTemplateGalleryModel(inputList *[]workbooktemplatesapis.WorkbookTemplateGallery) []WorkbookTemplateGalleryModel {
	if inputList == nil {
		return []WorkbookTemplateGalleryModel{}
	}

	outputList := make([]WorkbookTemplateGalleryModel, 0, len(*inputList))
	for _, input := range *inputList {
		output := WorkbookTemplateGalleryModel{}

		if input.Category != nil {
			output.Category = *input.Category
		}

		if input.Name != nil {
			output.Name = *input.Name
		}

		if input.Order != nil {
			output.Order = *input.Order
		}

		if input.ResourceType != nil {
			output.ResourceType = *input.ResourceType
		}

		if input.Type != nil {
			output.Type = *input.Type
		}

		outputList = append(outputList, output)
	}

	return outputList
}
