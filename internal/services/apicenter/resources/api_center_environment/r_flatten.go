package api_center_environment

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/environments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/services"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apicenter/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *environments.EnvironmentId, model *environments.Environment) error {
	state := ApiCenterEnvironmentResourceModel{
		ServiceId: services.NewServiceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName).ID(),
		Name:      id.EnvironmentName,
	}
	if model != nil {
		if props := model.Properties; props != nil {
			customProperties, err := helpers.FlattenCustomProperties(props.CustomProperties)
			if err != nil {
				return err
			}
			state.CustomProperties = customProperties

			state.Title = props.Title
			state.Description = pointer.From(props.Description)
			state.Type = string(props.Kind)
			if server := props.Server; server != nil {
				if pointer.From(server.ManagementPortalUri) != nil && len(pointer.From(server.ManagementPortalUri)) != 0 {
					state.MgmtPortalUri = pointer.From(server.ManagementPortalUri)[0]
				}
				state.ServerType = pointer.FromEnum(server.Type)
			}

			if onboarding := props.Onboarding; onboarding != nil {
				if pointer.From(onboarding.DeveloperPortalUri) != nil && len(pointer.From(onboarding.DeveloperPortalUri)) != 0 {
					state.DevPortalUri = pointer.From(onboarding.DeveloperPortalUri)[0]
				}
				state.Instructions = pointer.From(onboarding.Instructions)
			}
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}
