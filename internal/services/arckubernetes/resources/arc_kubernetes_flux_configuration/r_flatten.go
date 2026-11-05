// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arc_kubernetes_flux_configuration

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2025-04-01/fluxconfiguration"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/containers"
)

func flattenArcAzureBlobDefinitionModel(input *fluxconfiguration.AzureBlobDefinition, azureBlob []AzureBlobDefinitionModel, storageDomainSuffix string) ([]AzureBlobDefinitionModel, error) {
	outputList := make([]AzureBlobDefinitionModel, 0)
	if input == nil {
		return outputList, nil
	}

	accountId, err := accounts.ParseAccountID(pointer.From(input.Url), storageDomainSuffix)
	if err != nil {
		return nil, fmt.Errorf("parsing account %q: %+v", pointer.From(input.Url), err)
	}

	id := containers.NewContainerID(*accountId, pointer.From(input.ContainerName))

	output := AzureBlobDefinitionModel{
		ContainerID:           id.ID(),
		LocalAuthRef:          pointer.From(input.LocalAuthRef),
		SyncIntervalInSeconds: pointer.From(input.SyncIntervalInSeconds),
		TimeoutInSeconds:      pointer.From(input.TimeoutInSeconds),
	}

	var servicePrincipal []ServicePrincipalDefinitionModel
	if len(azureBlob) > 0 {
		output.AccountKey = azureBlob[0].AccountKey
		output.SasToken = azureBlob[0].SasToken
		servicePrincipal = azureBlob[0].ServicePrincipal
	}

	output.ServicePrincipal = flattenServicePrincipalDefinitionModel(input.ServicePrincipal, servicePrincipal)

	return append(outputList, output), nil
}

func flattenKustomizationDefinitionModel(inputList *map[string]fluxconfiguration.KustomizationDefinition) []KustomizationDefinitionModel {
	outputList := make([]KustomizationDefinitionModel, 0)
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := KustomizationDefinitionModel{
			DependsOn:              pointer.From(input.DependsOn),
			Force:                  pointer.From(input.Force),
			Name:                   pointer.From(input.Name),
			Path:                   pointer.From(input.Path),
			Prune:                  pointer.From(input.Prune),
			RetryIntervalInSeconds: pointer.From(input.RetryIntervalInSeconds),
			SyncIntervalInSeconds:  pointer.From(input.SyncIntervalInSeconds),
			TimeoutInSeconds:       pointer.From(input.TimeoutInSeconds),
		}

		outputList = append(outputList, output)
	}

	return outputList
}

func flattenServicePrincipalDefinitionModel(input *fluxconfiguration.ServicePrincipalDefinition, servicePrincipal []ServicePrincipalDefinitionModel) []ServicePrincipalDefinitionModel {
	outputList := make([]ServicePrincipalDefinitionModel, 0)
	if input == nil {
		return outputList
	}
	output := ServicePrincipalDefinitionModel{
		ClientCertificateSendChain: pointer.From(input.ClientCertificateSendChain),
		ClientId:                   pointer.From(input.ClientId),
		TenantId:                   pointer.From(input.TenantId),
	}

	if len(servicePrincipal) > 0 {
		output.ClientCertificate = servicePrincipal[0].ClientCertificate
		output.ClientCertificatePassword = servicePrincipal[0].ClientCertificatePassword
		output.ClientSecret = servicePrincipal[0].ClientSecret
	}

	return append(outputList, output)
}

func flattenBucketDefinitionModel(input *fluxconfiguration.BucketDefinition, bucket []BucketDefinitionModel) []BucketDefinitionModel {
	outputList := make([]BucketDefinitionModel, 0)
	if input == nil {
		return outputList
	}

	output := BucketDefinitionModel{
		AccessKey:             pointer.From(input.AccessKey),
		BucketName:            pointer.From(input.BucketName),
		TlsEnabled:            !pointer.From(input.Insecure),
		LocalAuthRef:          pointer.From(input.LocalAuthRef),
		SyncIntervalInSeconds: pointer.From(input.SyncIntervalInSeconds),
		TimeoutInSeconds:      pointer.From(input.TimeoutInSeconds),
		Url:                   pointer.From(input.Url),
	}

	if len(bucket) > 0 {
		output.SecretKey = bucket[0].SecretKey
	}

	return append(outputList, output)
}

func flattenGitRepositoryDefinitionModel(input *fluxconfiguration.GitRepositoryDefinition, gitRepository []GitRepositoryDefinitionModel) ([]GitRepositoryDefinitionModel, error) {
	outputList := make([]GitRepositoryDefinitionModel, 0)
	if input == nil {
		return outputList, nil
	}

	output := GitRepositoryDefinitionModel{
		HttpsCACert:           pointer.From(input.HTTPSCACert),
		HttpsUser:             pointer.From(input.HTTPSUser),
		LocalAuthRef:          pointer.From(input.LocalAuthRef),
		SshKnownHosts:         pointer.From(input.SshKnownHosts),
		SyncIntervalInSeconds: pointer.From(input.SyncIntervalInSeconds),
		TimeoutInSeconds:      pointer.From(input.TimeoutInSeconds),
		Url:                   pointer.From(input.Url),
	}

	referenceType, referenceValue, err := flattenRepositoryRefDefinitionModel(input.RepositoryRef)
	if err != nil {
		return nil, err
	}

	output.ReferenceType = referenceType
	output.ReferenceValue = referenceValue

	if len(gitRepository) > 0 {
		output.HttpsKey = gitRepository[0].HttpsKey
		output.SshPrivateKey = gitRepository[0].SshPrivateKey
	}

	return append(outputList, output), nil
}

func flattenRepositoryRefDefinitionModel(input *fluxconfiguration.RepositoryRefDefinition) (string, string, error) {
	if input == nil {
		return "", "", nil
	}

	var referenceType string
	var referenceValue string

	switch {
	case input.Branch != nil:
		referenceType = FluxGitBranch
		referenceValue = *input.Branch
	case input.Commit != nil:
		referenceType = FluxGitCommit
		referenceValue = *input.Commit
	case input.Semver != nil:
		referenceType = FluxGitSemverRange
		referenceValue = *input.Semver
	case input.Tag != nil:
		referenceType = FluxGitReferenceTag
		referenceValue = *input.Tag
	default:
		return "", "", fmt.Errorf("failed to retrieve git reference")
	}

	return referenceType, referenceValue, nil
}
