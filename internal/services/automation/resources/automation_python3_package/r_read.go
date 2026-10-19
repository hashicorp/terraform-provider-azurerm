// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_python3_package

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/python3package"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (m Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := python3package.ParsePython3PackageID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			client := meta.Client.Automation.Python3Package
			result, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(result.HttpResponse) {
					return meta.MarkAsGone(id)
				}
				return err
			}

			if result.Model == nil {
				return fmt.Errorf("retrieving %s got nil model", id)
			}

			var stateModel Python3PackageModel
			if err = meta.Decode(&stateModel); err != nil {
				return err
			}

			output := Python3PackageModel{
				ResourceGroupName:     id.ResourceGroupName,
				AutomationAccountName: id.AutomationAccountName,
				Name:                  id.Python3PackageName,

				// the fields below don't return by the API, remove it when issue fixed
				// https://github.com/Azure/azure-rest-api-specs/issues/25538
				ContentVersion: stateModel.ContentVersion,
				ContentUri:     stateModel.ContentUri,
				HashValue:      stateModel.HashValue,
				HashAlgorithm:  stateModel.HashAlgorithm,
			}

			model := result.Model
			if model.Properties != nil {
				if content := model.Properties.ContentLink; content != nil {
					output.ContentUri = pointer.From(content.Uri)
					output.ContentVersion = pointer.From(content.Version)
					if hash := content.ContentHash; hash != nil {
						output.HashAlgorithm = hash.Algorithm
						output.HashValue = hash.Value
					}
				}
				output.Tags = pointer.From(model.Tags)
			}

			return meta.Encode(&output)
		},
	}
}
