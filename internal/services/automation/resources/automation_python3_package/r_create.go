// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_python3_package

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/python3package"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (m Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automation.Python3Package

			var model Python3PackageModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			subscriptionID := metadata.Client.Account.SubscriptionId
			id := python3package.NewPython3PackageID(subscriptionID, model.ResourceGroupName, model.AutomationAccountName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if !response.WasNotFound(existing.HttpResponse) {
					if err != nil {
						return fmt.Errorf("retrieving %s: %v", id, err)
					}
					return metadata.ResourceRequiresImport(m.ResourceType(), id)
				}
			}

			req := python3package.PythonPackageCreateParameters{}
			req.Properties = python3package.PythonPackageCreateProperties{
				ContentLink: python3package.ContentLink{
					Uri: &model.ContentUri,
				},
			}

			if model.ContentVersion != "" {
				req.Properties.ContentLink.Version = &model.ContentVersion
			}

			if model.HashAlgorithm != "" {
				req.Properties.ContentLink.ContentHash = &python3package.ContentHash{
					Algorithm: model.HashAlgorithm,
					Value:     model.HashValue,
				}
			}
			req.Tags = &model.Tags

			if _, err := client.CreateOrUpdate(ctx, id, req); err != nil {
				return fmt.Errorf("creating %s: %v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
