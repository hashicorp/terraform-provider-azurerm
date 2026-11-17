// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_api_definition

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/apidefinitions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/apiversions"
	sdkclient "github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, ctx context.Context, id *apidefinitions.DefinitionId, model *apidefinitions.ApiDefinition, includeResources bool) error {
	client := metadata.Client.ApiCenter.ApiDefinitionsClient

	state := ApiCenterApiDefinitionModel{
		Name:         id.DefinitionName,
		ApiVersionId: apiversions.NewVersionID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceName, id.ApiName, id.VersionName).ID(),
	}
	exportResp := apidefinitions.ExportSpecificationOperationResponse{}
	if includeResources {
		var err error
		exportResp, err = client.ExportSpecification(ctx, *id)
		if err != nil && !response.WasNotFound(exportResp.HttpResponse) {
			return err
		}
		// Because the ExportSpecification does not unmarshall the body we have to do it here.
		if err == nil && exportResp.Model == nil &&
			exportResp.HttpResponse != nil &&
			exportResp.HttpResponse.StatusCode == http.StatusOK {
			var exported apidefinitions.ApiSpecExportResult
			resp := sdkclient.Response{Response: exportResp.HttpResponse}
			if err := resp.Unmarshal(&exported); err != nil {
				return fmt.Errorf("unmarshalling exported specification for %s: %+v", id, err)
			}
			exportResp.Model = &exported
		}
	}

	if props := model.Properties; props != nil {
		state.Description = pointer.From(props.Description)
		state.Title = props.Title
		state.Specification = flattenApiDefinitionPropertiesSpecification(&exportResp, props.Specification)
		if len(state.Specification) > 0 {
			state.Specification[0].Url = metadata.ResourceData.Get("specification.0.url").(string)
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func flattenApiDefinitionPropertiesSpecification(input *apidefinitions.ExportSpecificationOperationResponse, spec *apidefinitions.ApiDefinitionPropertiesSpecification) []ApiDefinitionPropertiesSpecification {
	if input == nil || spec == nil {
		return []ApiDefinitionPropertiesSpecification{}
	}

	content := ""
	url := ""
	if input.Model != nil {
		if pointer.From(input.Model.Format) == apidefinitions.ApiSpecExportResultFormatInline {
			content = pointer.From(input.Model.Value)
		}

		if pointer.From(input.Model.Format) == apidefinitions.ApiSpecExportResultFormatLink {
			url = pointer.From(input.Model.Value)
		}
	}

	return []ApiDefinitionPropertiesSpecification{
		{
			Name:    pointer.From(spec.Name),
			Version: pointer.From(spec.Version),
			Url:     url,
			Content: content,
		},
	}
}
