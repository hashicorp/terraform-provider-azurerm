// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_runtime_environment_package

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/packageresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automation.PackageResource

			id, err := packageresource.ParsePackageID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			var stateModel AutomationRuntimeEnvironmentPackageModel
			if err = metadata.Decode(&stateModel); err != nil {
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

			if props := resp.Model.Properties; props != nil {
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

			return metadata.Encode(&output)
		},
	}
}
