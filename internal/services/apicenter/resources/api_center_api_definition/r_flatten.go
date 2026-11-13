package api_center_api_definition

import (
	"context"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/apidefinitions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/apiversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, ctx context.Context, id *apidefinitions.DefinitionId, model *apidefinitions.ApiDefinition) error {
	client := metadata.Client.ApiCenter.ApiDefinitionsClient

	state := ApiCenterApiDefinitionModel{
		Name:         id.DefinitionName,
		ApiVersionId: apiversions.NewVersionID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, "default", id.ApiName, id.VersionName).ID(),
	}

	exportResp, err := client.ExportSpecification(ctx, *id)
	if err != nil {
		return err
	}

	if props := model.Properties; props != nil {
		specification := flattenApiDefinitionPropertiesSpecification(&exportResp, props.Specification)
		state.Description = pointer.From(props.Description)
		state.Specification = specification
		state.Title = props.Title
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
		if input.Model.Format == pointer.To(apidefinitions.ApiSpecExportResultFormatInline) {
			content = pointer.From(input.Model.Value)
		}

		if input.Model.Format == pointer.To(apidefinitions.ApiSpecExportResultFormatLink) {
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
