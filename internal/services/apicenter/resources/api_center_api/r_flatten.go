// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_api

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/apis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/services"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apicenter/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *apis.ApiId, model *apis.Api) error {
	state := ApiCenterApiModel{
		Name:      id.ApiName,
		ServiceId: services.NewServiceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName).ID(),
	}

	if props := model.Properties; props != nil {
		customProperties, err := helpers.FlattenCustomProperties(props.CustomProperties)
		if err != nil {
			return err
		}
		state.CustomProperties = customProperties
		state.Contacts = flattenContactList(props.Contacts)
		state.Description = pointer.From(props.Description)
		state.ExternalDocumentation = flattenExternalDocumentationList(props.ExternalDocumentation)
		state.Kind = string(props.Kind)
		state.License = flattenLicense(props.License)
		state.LifecycleStage = pointer.FromEnum(props.LifecycleStage)
		state.Summary = pointer.From(props.Summary)
		state.Title = props.Title
		if props.TermsOfService != nil {
			state.TermsOfServiceUrl = props.TermsOfService.Url
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func flattenContactList(input *[]apis.Contact) []Contact {
	if input == nil {
		return []Contact{}
	}
	result := make([]Contact, 0, len(*input))
	for _, v := range *input {
		result = append(result, Contact{
			Email: pointer.From(v.Email),
			Name:  pointer.From(v.Name),
			Url:   pointer.From(v.Url),
		})
	}
	return result
}

func flattenExternalDocumentationList(input *[]apis.ExternalDocumentation) []ExternalDocumentation {
	if input == nil {
		return []ExternalDocumentation{}
	}
	result := make([]ExternalDocumentation, 0, len(*input))

	for _, v := range *input {
		result = append(result, ExternalDocumentation{
			Description: pointer.From(v.Description),
			Title:       pointer.From(v.Title),
			Url:         v.Url,
		})
	}
	return result
}

func flattenLicense(input *apis.License) []License {
	if input == nil || (pointer.From(input.Identifier) == "" && pointer.From(input.Name) == "" && pointer.From(input.Url) == "") {
		return []License{}
	}
	return []License{
		{
			Identifier: pointer.From(input.Identifier),
			Name:       pointer.From(input.Name),
			Url:        pointer.From(input.Url),
		},
	}
}
