// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arc_kubernetes_flux_configuration

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2025-04-01/fluxconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ArcKubernetes.FluxConfigurationClient

			id, err := fluxconfiguration.ParseScopedFluxConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ArcKubernetesFluxConfigurationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if err := validateArcKubernetesFluxConfigurationModel(&model); err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			properties.Properties.ConfigurationProtectedSettings = nil
			if metadata.ResourceData.HasChange("blob_storage") {
				azureBlob, err := expandArcAzureBlobDefinitionModel(model.BlobStorage, metadata.Client.Storage.StorageDomainSuffix)
				if err != nil {
					return fmt.Errorf("expanding `blob_storage`: %+v", err)
				}

				properties.Properties.AzureBlob = azureBlob
				if properties.Properties.AzureBlob != nil {
					properties.Properties.SourceKind = pointer.To(fluxconfiguration.SourceKindTypeAzureBlob)
				}
			}

			if metadata.ResourceData.HasChange("bucket") {
				bucketValue, configurationProtectedSettings := expandBucketDefinitionModel(model.Bucket)
				properties.Properties.Bucket = bucketValue
				if properties.Properties.Bucket != nil {
					properties.Properties.SourceKind = pointer.To(fluxconfiguration.SourceKindTypeBucket)
					properties.Properties.ConfigurationProtectedSettings = configurationProtectedSettings
				}
			}

			if metadata.ResourceData.HasChange("git_repository") {
				gitRepositoryValue, configurationProtectedSettings, err := expandGitRepositoryDefinitionModel(model.GitRepository)
				if err != nil {
					return err
				}

				properties.Properties.GitRepository = gitRepositoryValue
				if properties.Properties.GitRepository != nil {
					properties.Properties.SourceKind = pointer.To(fluxconfiguration.SourceKindTypeGitRepository)
					properties.Properties.ConfigurationProtectedSettings = configurationProtectedSettings
				}
			}

			if metadata.ResourceData.HasChange("kustomizations") {
				properties.Properties.Kustomizations = expandKustomizationDefinitionModel(model.Kustomizations)
			}

			if metadata.ResourceData.HasChange("continuous_reconciliation_enabled") {
				properties.Properties.Suspend = pointer.To(!model.ContinuousReconciliationEnabled)
			}

			if properties.Properties.ConfigurationProtectedSettings == nil {
				if err := setConfigurationProtectedSettings(metadata, model, properties); err != nil {
					return err
				}
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
