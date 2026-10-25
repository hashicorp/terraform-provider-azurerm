// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_flux_configuration

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2025-04-01/fluxconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/containers"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model KubernetesFluxConfigurationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if err := validateKubernetesFluxConfigurationModel(&model); err != nil {
				return err
			}

			client := metadata.Client.Containers.KubernetesFluxConfigurationClient
			clusterID, err := commonids.ParseKubernetesClusterID(model.ClusterID)
			if err != nil {
				return err
			}

			// defined as strings because they're not enums in the swagger https://github.com/Azure/azure-rest-api-specs/pull/23545
			id := fluxconfiguration.NewScopedFluxConfigurationID(clusterID.ID(), model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			properties := &fluxconfiguration.FluxConfiguration{
				Properties: &fluxconfiguration.FluxConfigurationProperties{
					Kustomizations: expandKustomizationDefinitionModel(model.Kustomizations),
					Scope:          pointer.ToEnum[fluxconfiguration.ScopeType](model.Scope),
					Suspend:        pointer.To(!model.ContinuousReconciliationEnabled),
				},
			}

			if _, exists := metadata.ResourceData.GetOk("git_repository"); exists {
				gitRepositoryValue, configurationProtectedSettings, err := expandGitRepositoryDefinitionModel(model.GitRepository)
				if err != nil {
					return err
				}

				properties.Properties.SourceKind = pointer.To(fluxconfiguration.SourceKindTypeGitRepository)
				properties.Properties.GitRepository = gitRepositoryValue
				properties.Properties.ConfigurationProtectedSettings = configurationProtectedSettings
			} else if _, exists = metadata.ResourceData.GetOk("bucket"); exists {
				properties.Properties.SourceKind = pointer.To(fluxconfiguration.SourceKindTypeBucket)
				properties.Properties.Bucket, properties.Properties.ConfigurationProtectedSettings = expandBucketDefinitionModel(model.Bucket)
			} else if _, exists = metadata.ResourceData.GetOk("blob_storage"); exists {
				properties.Properties.SourceKind = pointer.To(fluxconfiguration.SourceKindTypeAzureBlob)
				azureBlob, err := expandAzureBlobDefinitionModel(model.BlobStorage, metadata.Client.Storage.StorageDomainSuffix)
				if err != nil {
					return fmt.Errorf("expanding `blob_storage`: %+v", err)
				}

				properties.Properties.AzureBlob = azureBlob
			}

			if model.Namespace != "" {
				properties.Properties.Namespace = &model.Namespace
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, *properties, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func expandAzureBlobDefinitionModel(inputList []AzureBlobDefinitionModel, storageDomainSuffix string) (*fluxconfiguration.AzureBlobDefinition, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	input := &inputList[0]
	output := fluxconfiguration.AzureBlobDefinition{
		ManagedIdentity:       expandManagedIdentityDefinitionModel(input.ManagedIdentity),
		ServicePrincipal:      expandServicePrincipalDefinitionModel(input.ServicePrincipal),
		SyncIntervalInSeconds: &input.SyncIntervalInSeconds,
		TimeoutInSeconds:      &input.TimeoutInSeconds,
	}
	if input.AccountKey != "" {
		output.AccountKey = &input.AccountKey
	}

	if input.ContainerID != "" {
		id, err := containers.ParseContainerID(input.ContainerID, storageDomainSuffix)
		if err != nil {
			return nil, err
		}

		output.ContainerName = &id.ContainerName
		output.Url = pointer.To(id.AccountId.ID())
	}

	if input.LocalAuthRef != "" {
		output.LocalAuthRef = &input.LocalAuthRef
	}

	if input.SasToken != "" {
		output.SasToken = &input.SasToken
	}

	return &output, nil
}

func expandManagedIdentityDefinitionModel(inputList []ManagedIdentityDefinitionModel) *fluxconfiguration.ManagedIdentityDefinition {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := fluxconfiguration.ManagedIdentityDefinition{}
	if input.ClientId != "" {
		output.ClientId = &input.ClientId
	}

	return &output
}

func expandKustomizationDefinitionModel(inputList []KustomizationDefinitionModel) *map[string]fluxconfiguration.KustomizationDefinition {
	if len(inputList) == 0 {
		return nil
	}

	outputList := make(map[string]fluxconfiguration.KustomizationDefinition)
	for _, v := range inputList {
		input := v
		// updated item in a set is considered a new item, and the old item still exists in the set but with empty values, so we need to skip it
		if input.Name == "" {
			continue
		}
		output := fluxconfiguration.KustomizationDefinition{
			DependsOn:              &input.DependsOn,
			Force:                  &input.Force,
			Name:                   &input.Name,
			PostBuild:              expandPostBuildDefinitionModel(input.PostBuild),
			Prune:                  &input.Prune,
			RetryIntervalInSeconds: &input.RetryIntervalInSeconds,
			SyncIntervalInSeconds:  &input.SyncIntervalInSeconds,
			TimeoutInSeconds:       &input.TimeoutInSeconds,
			Wait:                   &input.Wait,
		}

		if input.Path != "" {
			output.Path = pointer.To(input.Path)
		}

		outputList[input.Name] = output
	}

	return &outputList
}

func expandPostBuildDefinitionModel(inputList []PostBuildDefinitionModel) *fluxconfiguration.PostBuildDefinition {
	if len(inputList) == 0 {
		return nil
	}

	input := inputList[0]

	output := fluxconfiguration.PostBuildDefinition{}

	if len(input.Substitute) > 0 {
		output.Substitute = &input.Substitute
	}

	if len(input.SubstituteFrom) > 0 {
		output.SubstituteFrom = expandSubstituteFromDefinitionModel(input.SubstituteFrom)
	}

	return &output
}

func expandSubstituteFromDefinitionModel(inputList []SubstituteFromDefinitionModel) *[]fluxconfiguration.SubstituteFromDefinition {
	if len(inputList) == 0 {
		return nil
	}

	input := inputList
	output := make([]fluxconfiguration.SubstituteFromDefinition, 0)

	for _, v := range input {
		output = append(output, fluxconfiguration.SubstituteFromDefinition{
			Kind:     &v.Kind,
			Name:     &v.Name,
			Optional: &v.Optional,
		})
	}

	return &output
}

func expandServicePrincipalDefinitionModel(inputList []ServicePrincipalDefinitionModel) *fluxconfiguration.ServicePrincipalDefinition {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := fluxconfiguration.ServicePrincipalDefinition{
		ClientCertificateSendChain: &input.ClientCertificateSendChain,
	}

	if input.ClientCertificate != "" {
		output.ClientCertificate = &input.ClientCertificate
	}

	if input.ClientCertificatePassword != "" {
		output.ClientCertificatePassword = &input.ClientCertificatePassword
	}

	if input.ClientId != "" {
		output.ClientId = &input.ClientId
	}

	if input.ClientSecret != "" {
		output.ClientSecret = &input.ClientSecret
	}

	if input.TenantId != "" {
		output.TenantId = &input.TenantId
	}

	return &output
}

func expandBucketDefinitionModel(inputList []BucketDefinitionModel) (*fluxconfiguration.BucketDefinition, *map[string]string) {
	if len(inputList) == 0 {
		return nil, nil
	}

	input := &inputList[0]
	output := fluxconfiguration.BucketDefinition{
		Insecure:              pointer.To(!input.TlsEnabled),
		SyncIntervalInSeconds: &input.SyncIntervalInSeconds,
		TimeoutInSeconds:      &input.TimeoutInSeconds,
	}

	if input.AccessKey != "" {
		output.AccessKey = &input.AccessKey
	}

	if input.BucketName != "" {
		output.BucketName = &input.BucketName
	}

	if input.LocalAuthRef != "" {
		output.LocalAuthRef = &input.LocalAuthRef
	}

	if input.Url != "" {
		output.Url = &input.Url
	}

	configSettings := make(map[string]string)
	if input.SecretKey != "" {
		configSettings["bucketSecretKey"] = input.SecretKey
	}

	var outputConfigSettings *map[string]string = nil
	if len(configSettings) > 0 {
		outputConfigSettings = &configSettings
	}

	return &output, outputConfigSettings
}

func expandGitRepositoryDefinitionModel(inputList []GitRepositoryDefinitionModel) (*fluxconfiguration.GitRepositoryDefinition, *map[string]string, error) {
	if len(inputList) == 0 {
		return nil, nil, nil
	}

	input := &inputList[0]
	output := fluxconfiguration.GitRepositoryDefinition{
		SyncIntervalInSeconds: &input.SyncIntervalInSeconds,
		TimeoutInSeconds:      &input.TimeoutInSeconds,
	}

	if input.HttpsCACert != "" {
		output.HTTPSCACert = &input.HttpsCACert
	}

	if input.HttpsUser != "" {
		output.HTTPSUser = &input.HttpsUser
	}

	if input.LocalAuthRef != "" {
		output.LocalAuthRef = &input.LocalAuthRef
	}

	repositoryRefValue, err := expandRepositoryRefDefinitionModel(input.ReferenceType, input.ReferenceValue)
	if err != nil {
		return nil, nil, err
	}

	output.RepositoryRef = repositoryRefValue

	if input.SshKnownHosts != "" {
		output.SshKnownHosts = &input.SshKnownHosts
	}

	if input.Url != "" {
		output.Url = &input.Url
	}

	if input.Provider != "" {
		output.Provider = pointer.ToEnum[fluxconfiguration.ProviderType](input.Provider)
	}

	configSettings := make(map[string]string)
	if input.HttpsKey != "" {
		configSettings["httpsKey"] = input.HttpsKey
	}

	if input.SshPrivateKey != "" {
		configSettings["sshPrivateKey"] = input.SshPrivateKey
	}

	return &output, &configSettings, nil
}

func expandRepositoryRefDefinitionModel(referenceType string, referenceValue string) (*fluxconfiguration.RepositoryRefDefinition, error) {
	output := fluxconfiguration.RepositoryRefDefinition{}

	switch referenceType {
	case FluxGitBranch:
		output.Branch = &referenceValue
	case FluxGitCommit:
		output.Commit = &referenceValue
	case FluxGitSemverRange:
		output.Semver = &referenceValue
	case FluxGitReferenceTag:
		output.Tag = &referenceValue
	default:
		return &output, fmt.Errorf("reference type %s not defined", referenceType)
	}

	return &output, nil
}

func setConfigurationProtectedSettings(metadata sdk.ResourceMetaData, model KubernetesFluxConfigurationModel, properties *fluxconfiguration.FluxConfiguration) error {
	if _, exists := metadata.ResourceData.GetOk("git_repository"); exists {
		_, configurationProtectedSettings, err := expandGitRepositoryDefinitionModel(model.GitRepository)
		if err != nil {
			return err
		}
		properties.Properties.ConfigurationProtectedSettings = configurationProtectedSettings
	} else if _, exists = metadata.ResourceData.GetOk("bucket"); exists {
		_, properties.Properties.ConfigurationProtectedSettings = expandBucketDefinitionModel(model.Bucket)
	}
	return nil
}
