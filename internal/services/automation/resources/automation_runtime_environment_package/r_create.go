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
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/custompollers"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automation.PackageResource

			var model AutomationRuntimeEnvironmentPackageModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			envId, err := packageresource.ParseRuntimeEnvironmentID(model.AutomationRuntimeEnvironmentId)
			if err != nil {
				return err
			}

			id := packageresource.NewPackageID(envId.SubscriptionId, envId.ResourceGroupName, envId.AutomationAccountName, envId.RuntimeEnvironmentName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			parameters := packageresource.PackageCreateOrUpdateParameters{
				Properties: packageresource.PackageCreateOrUpdateProperties{
					ContentLink: packageresource.ContentLink{
						Uri: pointer.To(model.ContentUri),
					},
				},
			}

			if model.ContentVersion != "" {
				parameters.Properties.ContentLink.Version = pointer.To(model.ContentVersion)
			}

			if model.HashAlgorithm != "" {
				parameters.Properties.ContentLink.ContentHash = &packageresource.ContentHash{
					Algorithm: model.HashAlgorithm,
					Value:     model.HashValue,
				}
			}

			if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			// custom poller is required until https://github.com/Azure/azure-rest-api-specs/issues/41641 is resolved
			pollerType := custompollers.NewAutomationRuntimeEnvironmentPackagePoller(client, id)
			poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
			if err := poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for %s to finish provisioning: %+v", id, err)
			}

			return nil
		},
	}
}
