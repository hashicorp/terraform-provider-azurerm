// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package image_builder_template

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/imagebuilder/2024-02-01/virtualmachineimagetemplate"
)

func flattenBasicImageTemplateSource(input virtualmachineimagetemplate.ImageTemplateSource) (string, []ImageBuilderTemplateSourcePlatformImage, string) {
	if input != nil {
		switch source := input.(type) {
		case virtualmachineimagetemplate.ImageTemplateManagedImageSource:
			return flattenSourceManagedImage(&source), make([]ImageBuilderTemplateSourcePlatformImage, 0), ""
		case virtualmachineimagetemplate.ImageTemplatePlatformImageSource:
			return "", flattenSourcePlatformImage(&source), ""
		case virtualmachineimagetemplate.ImageTemplateSharedImageVersionSource:
			return "", make([]ImageBuilderTemplateSourcePlatformImage, 0), flattenSourceSharedImageVersion(&source)
		}
	}

	return "", make([]ImageBuilderTemplateSourcePlatformImage, 0), ""
}

func flattenSourceManagedImage(input *virtualmachineimagetemplate.ImageTemplateManagedImageSource) string {
	if input != nil && input.ImageId != "" {
		return input.ImageId
	}

	return ""
}

func flattenSourcePlatformImage(input *virtualmachineimagetemplate.ImageTemplatePlatformImageSource) []ImageBuilderTemplateSourcePlatformImage {
	if input == nil {
		return make([]ImageBuilderTemplateSourcePlatformImage, 0)
	}

	result := make([]ImageBuilderTemplateSourcePlatformImage, 0)

	result = append(result, ImageBuilderTemplateSourcePlatformImage{
		Plan: flattenImageBuilderTemplateSourcePlatformImagePlan(input.PlanInfo),
	})

	if input.Publisher != nil {
		result[0].Publisher = *input.Publisher
	}

	if input.Offer != nil {
		result[0].Offer = *input.Offer
	}

	if input.Sku != nil {
		result[0].Sku = *input.Sku
	}

	if input.Version != nil {
		result[0].Version = *input.Version
	}
	return result
}

func flattenImageBuilderTemplateSourcePlatformImagePlan(input *virtualmachineimagetemplate.PlatformImagePurchasePlan) []ImageBuilderTemplateSourcePlatformImagePlan {
	if input == nil {
		return make([]ImageBuilderTemplateSourcePlatformImagePlan, 0)
	}

	result := make([]ImageBuilderTemplateSourcePlatformImagePlan, 0)

	result = append(result, ImageBuilderTemplateSourcePlatformImagePlan{
		Name:      input.PlanName,
		Product:   input.PlanProduct,
		Publisher: input.PlanPublisher,
	})
	return result
}

func flattenSourceSharedImageVersion(input *virtualmachineimagetemplate.ImageTemplateSharedImageVersionSource) string {
	if input != nil && input.ImageVersionId != "" {
		return input.ImageVersionId
	}

	return ""
}

func flattenBasicImageTemplateDistributor(input *[]virtualmachineimagetemplate.ImageTemplateDistributor) ([]ImageBuilderTemplateDistributions, error) {
	results := make([]ImageBuilderTemplateDistributions, 0)

	distributionManagedImages := make([]ImageBuilderTemplateDistributionsManagedImage, 0)
	distributionSharedImages := make([]ImageBuilderTemplateDistributionsSharedImage, 0)
	distributionVhds := make([]ImageBuilderTemplateDistributionsVhd, 0)

	if input != nil {
		for _, v := range *input {
			switch distribute := v.(type) {
			case virtualmachineimagetemplate.ImageTemplateManagedImageDistributor:
				flattenedManagedImg, err := flattenDistributionManagedImage(&distribute)
				if err != nil {
					return nil, fmt.Errorf("setting image template managed image source: %+v", err)
				}

				distributionManagedImages = append(distributionManagedImages, flattenedManagedImg)
			case virtualmachineimagetemplate.ImageTemplateSharedImageDistributor:
				distributionSharedImages = append(distributionSharedImages, flattenDistributionSharedImage(&distribute))
			case virtualmachineimagetemplate.ImageTemplateVhdDistributor:
				distributionVhds = append(distributionVhds, flattenDistributionVhd(&distribute))
			}
		}
	}

	output := ImageBuilderTemplateDistributions{
		ManagedImage: distributionManagedImages,
		SharedImage:  distributionSharedImages,
		Vhd:          distributionVhds,
	}

	results = append(results, output)
	return results, nil
}

func flattenImageTemplateSharedImageDistributorVersioning(input virtualmachineimagetemplate.DistributeVersioner) []ImageBuilderTemplateDistributionsSharedImageVersioning {
	results := make([]ImageBuilderTemplateDistributionsSharedImageVersioning, 0)
	output := ImageBuilderTemplateDistributionsSharedImageVersioning{}

	switch versioner := input.(type) {
	case virtualmachineimagetemplate.DistributeVersionerLatest:
		output.Scheme = versioner.Scheme
		output.Major = *versioner.Major

	case virtualmachineimagetemplate.DistributeVersionerSource:
		output.Scheme = versioner.Scheme
	}

	results = append(results, output)

	return results
}

func flattenImageTemplateSharedImageDistributorReplicaRegions(input *[]string) []ImageBuilderTemplateDistributionsSharedImageReplicaRegions {
	results := make([]ImageBuilderTemplateDistributionsSharedImageReplicaRegions, 0)

	if input != nil {
		for _, v := range *input {
			results = append(results, ImageBuilderTemplateDistributionsSharedImageReplicaRegions{
				Name: v,
			})
		}
	}

	return results
}

func flattenDistributionManagedImage(input *virtualmachineimagetemplate.ImageTemplateManagedImageDistributor) (ImageBuilderTemplateDistributionsManagedImage, error) {
	result := ImageBuilderTemplateDistributionsManagedImage{}
	if input == nil {
		return result, nil
	}

	imageName := ""
	resourceGroupName := ""
	if input.ImageId != "" {
		imageNameReturned, resourceGroupNameReturned, err := imageBuilderTemplateManagedImageNameAndResourceGroupName(input.ImageId)
		if err != nil {
			return result, err
		}

		imageName = imageNameReturned
		resourceGroupName = resourceGroupNameReturned
	}

	result.Name = imageName
	result.ResourceGroupName = resourceGroupName
	result.Location = input.Location
	result.RunOutputName = input.RunOutputName
	result.Tags = *input.ArtifactTags

	return result, nil
}

func flattenDistributionSharedImage(input *virtualmachineimagetemplate.ImageTemplateSharedImageDistributor) ImageBuilderTemplateDistributionsSharedImage {
	results := ImageBuilderTemplateDistributionsSharedImage{}
	if input == nil {
		return results
	}

	results.Id = input.GalleryImageId

	if input.ReplicationRegions != nil {
		results.ReplicaRegions = flattenImageTemplateSharedImageDistributorReplicaRegions(input.ReplicationRegions)
	}

	results.RunOutputName = input.RunOutputName

	if input.ExcludeFromLatest != nil {
		results.ExcludeFromLatest = *input.ExcludeFromLatest
	}

	if input.StorageAccountType != nil {
		results.StorageAccountType = string(*input.StorageAccountType)
	}

	if input.ArtifactTags != nil {
		results.Tags = *input.ArtifactTags
	}

	results.Versioning = flattenImageTemplateSharedImageDistributorVersioning(input.Versioning)
	return results
}

func flattenDistributionVhd(input *virtualmachineimagetemplate.ImageTemplateVhdDistributor) ImageBuilderTemplateDistributionsVhd {
	results := ImageBuilderTemplateDistributionsVhd{}
	if input == nil {
		return results
	}

	if input.RunOutputName != "" {
		results.RunOutputName = input.RunOutputName
	}

	if input.ArtifactTags != nil {
		results.Tags = *input.ArtifactTags
	}

	return results
}

func flattenBasicImageTemplateCustomizer(input *[]virtualmachineimagetemplate.ImageTemplateCustomizer) []ImageBuilderTemplateCustomizer {
	customizerList := make([]ImageBuilderTemplateCustomizer, 0)

	if input != nil {
		for _, v := range *input {
			switch customizer := v.(type) {
			case virtualmachineimagetemplate.ImageTemplateFileCustomizer:
				customizerList = append(customizerList, flattenCustomizerFile(&customizer))
			case virtualmachineimagetemplate.ImageTemplatePowerShellCustomizer:
				customizerList = append(customizerList, flattenCustomizerPowerShell(&customizer))
			case virtualmachineimagetemplate.ImageTemplateShellCustomizer:
				customizerList = append(customizerList, flattenCustomizerShell(&customizer))
			case virtualmachineimagetemplate.ImageTemplateRestartCustomizer:
				customizerList = append(customizerList, flattenCustomizerWindowsRestart(&customizer))
			case virtualmachineimagetemplate.ImageTemplateWindowsUpdateCustomizer:
				customizerList = append(customizerList, flattenCustomizerWindowsUpdate(&customizer))
			}
		}

		return customizerList
	}

	return customizerList
}

func flattenCustomizerFile(input *virtualmachineimagetemplate.ImageTemplateFileCustomizer) ImageBuilderTemplateCustomizer {
	result := ImageBuilderTemplateCustomizer{}

	if input == nil {
		return result
	}

	result.Type = "File"

	if input.Name != nil {
		result.Name = *input.Name
	}

	if input.SourceUri != nil {
		result.FileSourceUri = *input.SourceUri
	}

	if input.Sha256Checksum != nil {
		result.FileSha256Checksum = *input.Sha256Checksum
	}

	if input.Destination != nil {
		result.FileDestinationPath = *input.Destination
	}

	return result
}

func flattenCustomizerShell(input *virtualmachineimagetemplate.ImageTemplateShellCustomizer) ImageBuilderTemplateCustomizer {
	result := ImageBuilderTemplateCustomizer{}
	if input == nil {
		return result
	}

	result.Type = "Shell"

	if input.Name != nil {
		result.Name = *input.Name
	}

	if input.ScriptUri != nil {
		result.ShellScriptUri = *input.ScriptUri
	}

	if input.Sha256Checksum != nil {
		result.ShellSha256Checksum = *input.Sha256Checksum
	}

	if input.Inline != nil {
		result.ShellCommands = *input.Inline
	}

	return result
}

func flattenCustomizerPowerShell(input *virtualmachineimagetemplate.ImageTemplatePowerShellCustomizer) ImageBuilderTemplateCustomizer {
	result := ImageBuilderTemplateCustomizer{}
	if input == nil {
		return result
	}

	result.Type = "PowerShell"

	if input.Name != nil {
		result.Name = *input.Name
	}

	if input.ScriptUri != nil {
		result.PowershellScriptUri = *input.ScriptUri
	}

	if input.Sha256Checksum != nil {
		result.PowershellSha256Checksum = *input.Sha256Checksum
	}

	if input.Inline != nil {
		result.PowershellCommands = *input.Inline
	}

	if input.RunAsSystem != nil {
		result.PowershellRunAsSystem = pointer.From(input.RunAsSystem)
	}

	if input.RunElevated != nil {
		result.PowershellRunElevated = *input.RunElevated
	}

	if input.ValidExitCodes != nil {
		result.PowershellValidExitCodes = *input.ValidExitCodes
	}

	return result
}

func flattenCustomizerWindowsRestart(input *virtualmachineimagetemplate.ImageTemplateRestartCustomizer) ImageBuilderTemplateCustomizer {
	result := ImageBuilderTemplateCustomizer{}
	if input == nil {
		return result
	}

	result.Type = "WindowsRestart"

	if input.Name != nil {
		result.Name = *input.Name
	}

	if input.RestartCommand != nil {
		result.WindowsRestartCommand = *input.RestartCommand
	}

	if input.RestartCheckCommand != nil {
		result.WindowsRestartCheckCommand = *input.RestartCheckCommand
	}

	if input.RestartTimeout != nil {
		result.WindowsRestartTimeout = *input.RestartTimeout
	}

	return result
}

func flattenCustomizerWindowsUpdate(input *virtualmachineimagetemplate.ImageTemplateWindowsUpdateCustomizer) ImageBuilderTemplateCustomizer {
	result := ImageBuilderTemplateCustomizer{}
	if input == nil {
		return result
	}

	result.Type = "WindowsUpdate"

	if input.Name != nil {
		result.Name = *input.Name
	}

	if input.SearchCriteria != nil {
		result.WindowsUpdateSearchCriteria = *input.SearchCriteria
	}

	if input.Filters != nil {
		result.WindowsUpdateFilters = *input.Filters
	}

	if input.UpdateLimit != nil {
		result.WindowsUpdateLimit = *input.UpdateLimit
	}

	return result
}
