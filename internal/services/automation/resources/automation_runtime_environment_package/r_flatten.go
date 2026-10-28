package automation_runtime_environment_package

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/packageresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *packageresource.PackageId, model *packageresource.Package) error {
	var stateModel AutomationRuntimeEnvironmentPackageModel
	if err := metadata.Decode(&stateModel); err != nil {
		return fmt.Errorf("decoding: %+v", err)
	}

	runtimeEnvId := packageresource.NewRuntimeEnvironmentID(id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName, id.RuntimeEnvironmentName)

	output := AutomationRuntimeEnvironmentPackageModel{
		Name:                           id.PackageName,
		AutomationRuntimeEnvironmentId: runtimeEnvId.ID(),

		// the fields below don't return by the API, remove it when issue fixed
		// https://github.com/Azure/azure-rest-api-specs/issues/41604
		ContentVersion: stateModel.ContentVersion,
		ContentUri:     stateModel.ContentUri,
		HashValue:      stateModel.HashValue,
		HashAlgorithm:  stateModel.HashAlgorithm,
	}

	if model != nil {
		if props := model.Properties; props != nil {
			if content := props.ContentLink; content != nil {
				output.ContentUri = pointer.From(content.Uri)
				output.ContentVersion = pointer.From(content.Version)

				if hash := content.ContentHash; hash != nil {
					output.HashAlgorithm = hash.Algorithm
					output.HashValue = hash.Value
				}
			}

			output.SizeInBytes = pointer.From(props.SizeInBytes)
			output.Version = pointer.From(props.Version)
			output.Default = pointer.From(props.Default)
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&output)
}
