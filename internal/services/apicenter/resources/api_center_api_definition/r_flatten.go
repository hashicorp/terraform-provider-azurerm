package api_center_api_definition

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/apidefinitions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/apiversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *apidefinitions.DefinitionId, model *apidefinitions.ApiDefinition) error {
	state := ApiCenterApiDefinitionModel{
		Name:         id.DefinitionName,
		ApiVersionId: apiversions.NewVersionID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, "default", id.ApiName, id.VersionName).ID(),
	}

	if props := model.Properties; props != nil {
		state.Description = pointer.From(props.Description)
		state.Specification = flattenApiDefinitionPropertiesSpecification(props.Specification)
		state.Title = props.Title
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func flattenApiDefinitionPropertiesSpecification(input *apidefinitions.ApiDefinitionPropertiesSpecification) []ApiDefinitionPropertiesSpecification {
	if input == nil {
		return []ApiDefinitionPropertiesSpecification{}
	}
	return []ApiDefinitionPropertiesSpecification{
		{
			Name:    pointer.From(input.Name),
			Version: pointer.From(input.Version),
		},
	}
}
