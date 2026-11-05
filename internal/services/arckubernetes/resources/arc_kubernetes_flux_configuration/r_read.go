// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arc_kubernetes_flux_configuration

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridkubernetes/2024-01-01/connectedclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2025-04-01/fluxconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ArcKubernetes.FluxConfigurationClient

			id, err := fluxconfiguration.ParseScopedFluxConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var configModel ArcKubernetesFluxConfigurationModel
			if err := metadata.Decode(&configModel); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			connectedClusterId, err := connectedclusters.ParseConnectedClusterID(id.Scope)
			if err != nil {
				return fmt.Errorf("parsing %q as a Connected Cluster ID: %+v", id.Scope, err)
			}

			state := ArcKubernetesFluxConfigurationModel{
				Name:      id.FluxConfigurationName,
				ClusterID: connectedClusterId.ID(),
			}

			if model := resp.Model; model != nil {
				if properties := model.Properties; properties != nil {
					blobStorage, err := flattenArcAzureBlobDefinitionModel(properties.AzureBlob, configModel.BlobStorage, metadata.Client.Storage.StorageDomainSuffix)
					if err != nil {
						return fmt.Errorf("flattening `blob_storage`: %+v", err)
					}

					state.BlobStorage = blobStorage
					state.Bucket = flattenBucketDefinitionModel(properties.Bucket, configModel.Bucket)
					gitRepositoryValue, err := flattenGitRepositoryDefinitionModel(properties.GitRepository, configModel.GitRepository)
					if err != nil {
						return err
					}

					state.GitRepository = gitRepositoryValue
					state.Kustomizations = flattenKustomizationDefinitionModel(properties.Kustomizations)
					state.Namespace = pointer.From(properties.Namespace)
					state.Scope = pointer.FromEnum(properties.Scope)
					state.ContinuousReconciliationEnabled = !pointer.From(properties.Suspend)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

const (
	FluxGitBranch       string = "branch"
	FluxGitCommit       string = "commit"
	FluxGitReferenceTag string = "tag"
	FluxGitSemverRange  string = "semver"
)

func validateArcKubernetesFluxConfigurationModel(model *ArcKubernetesFluxConfigurationModel) error {
	allKeys := make(map[string]bool)
	for _, k := range model.Kustomizations {
		if _, exists := allKeys[k.Name]; exists {
			return fmt.Errorf("kustomization name `%s` is not unique", k.Name)
		}

		allKeys[k.Name] = true
	}

	return nil
}
