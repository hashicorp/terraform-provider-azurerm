// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package image_builder_template

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/imagebuilder/2024-02-01/virtualmachineimagetemplate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandBasicImageTemplateSource(managedImageId string, platformImage []ImageBuilderTemplateSourcePlatformImage, sharedImageId string) virtualmachineimagetemplate.ImageTemplateSource {
	if len(managedImageId) != 0 {
		return &virtualmachineimagetemplate.ImageTemplateManagedImageSource{
			ImageId: managedImageId,
		}
	}

	if len(platformImage) != 0 {
		v := platformImage[0]
		result := &virtualmachineimagetemplate.ImageTemplatePlatformImageSource{
			Publisher: &v.Publisher,
			Offer:     &v.Offer,
			Sku:       &v.Sku,
			Version:   &v.Version,
		}

		if len(v.Plan) > 0 {
			planTemp := v.Plan[0]
			result.PlanInfo = &virtualmachineimagetemplate.PlatformImagePurchasePlan{
				PlanName:      planTemp.Name,
				PlanProduct:   planTemp.Product,
				PlanPublisher: planTemp.Publisher,
			}
		}

		return result
	}

	if len(sharedImageId) != 0 {
		return &virtualmachineimagetemplate.ImageTemplateSharedImageVersionSource{
			ImageVersionId: sharedImageId,
		}
	}

	// since there is ExactlyOneOf on the source schema, below nil won't be reached
	return nil
}

func expandBasicImageTemplateDistributor(distributions []ImageBuilderTemplateDistributions, subscriptionId string) (*[]virtualmachineimagetemplate.ImageTemplateDistributor, error) {
	results := make([]virtualmachineimagetemplate.ImageTemplateDistributor, 0)
	runOutputNameSet := make(map[string]bool)

	if len(distributions) > 0 {
		for _, v := range distributions {
			managedImages := v.ManagedImage
			sharedImages := v.SharedImage
			vhds := v.Vhd

			if len(managedImages) == 0 && len(sharedImages) == 0 && len(vhds) == 0 {
				return &results, fmt.Errorf("at least one of `managed_image`, `shared_image` and `vhd` is required to specify in `distributions`")
			}

			if len(managedImages) > 0 {
				for _, managedImage := range managedImages {
					resourceGroupName := managedImage.ResourceGroupName
					runOutputName := managedImage.RunOutputName

					_, existing := runOutputNameSet[runOutputName]
					if existing {
						return &results, fmt.Errorf("`run_output_name` must be unique among all distribution destinations. %q already exists", runOutputName)
					} else {
						runOutputNameSet[runOutputName] = true

						results = append(results, virtualmachineimagetemplate.ImageTemplateManagedImageDistributor{
							ImageId:       "/subscriptions/" + subscriptionId + "/resourceGroups/" + resourceGroupName + "/providers/Microsoft.Compute/images/" + managedImage.Name,
							Location:      managedImage.Location,
							RunOutputName: runOutputName,
							ArtifactTags:  &managedImage.Tags,
						})
					}
				}
			}

			if len(sharedImages) > 0 {
				for _, sharedImage := range sharedImages {
					// sharedImage := sharedImageRaw.(map[string]any)
					runOutputName := sharedImage.RunOutputName

					_, existing := runOutputNameSet[runOutputName]
					if existing {
						return &results, fmt.Errorf("`run_output_name` must be unique among all distribution destinations. %q already exists", runOutputName)
					} else {
						runOutputNameSet[runOutputName] = true
						results = append(results, virtualmachineimagetemplate.ImageTemplateSharedImageDistributor{
							GalleryImageId:     sharedImage.Id,
							ReplicationRegions: expandImageTemplateSharedImageDistributorReplicaRegions(sharedImage.ReplicaRegions),
							RunOutputName:      sharedImage.RunOutputName,
							ExcludeFromLatest:  &sharedImage.ExcludeFromLatest,
							StorageAccountType: pointer.ToEnum[virtualmachineimagetemplate.SharedImageStorageAccountType](sharedImage.StorageAccountType),
							ArtifactTags:       &sharedImage.Tags,
							Versioning:         expandImageTemplateSharedImageDistributorVersioning(sharedImage.Versioning),
						})
					}
				}
			}

			if len(vhds) > 0 {
				for _, vhd := range vhds {
					// vhd := vhdRaw.(map[string]any)
					runOutputName := vhd.RunOutputName

					_, existing := runOutputNameSet[runOutputName]
					if existing {
						return &results, fmt.Errorf("`run_output_name` must be unique among all distribution destinations. %q already exists", runOutputName)
					} else {
						runOutputNameSet[runOutputName] = true
						results = append(results, virtualmachineimagetemplate.ImageTemplateVhdDistributor{
							RunOutputName: vhd.RunOutputName,
							ArtifactTags:  &vhd.Tags,
						})
					}
				}
			}
		}
	}

	return &results, nil
}

func expandImageTemplateSharedImageDistributorReplicaRegions(input []ImageBuilderTemplateDistributionsSharedImageReplicaRegions) *[]string {
	if input == nil {
		return nil
	}

	result := make([]string, 0)
	for _, v := range input {
		result = append(result, v.Name)
	}

	return &result
}

func expandImageTemplateSharedImageDistributorVersioning(input []ImageBuilderTemplateDistributionsSharedImageVersioning) virtualmachineimagetemplate.DistributeVersioner {
	var result virtualmachineimagetemplate.DistributeVersioner
	if len(input) > 0 {
		config := input[0]

		if v := config.Scheme; v == "Latest" {
			result = virtualmachineimagetemplate.DistributeVersionerLatest{
				Scheme: config.Scheme,
				Major:  &config.Major,
			}
		}

		if v := config.Scheme; v == "Source" {
			result = virtualmachineimagetemplate.DistributeVersionerSource{
				Scheme: config.Scheme,
			}
		}
	}
	return result
}

// Passing d as input rather its business subset because this function needs d.GetOK() to verify users' input to the `customizer` block is valid.
func expandBasicImageTemplateCustomizer(d *schema.ResourceData, input []ImageBuilderTemplateCustomizer) (*[]virtualmachineimagetemplate.ImageTemplateCustomizer, error) {
	if len(input) == 0 {
		return nil, nil
	}

	results := make([]virtualmachineimagetemplate.ImageTemplateCustomizer, 0)

	// This index serves identifying the invalid input customizer block from the customizer block list.
	i := 0

	for _, customizer := range input {
		t := customizer.Type

		switch t {
		case "File":
			if err := validateImageTemplateCustomizerInputForFileType(d, customizer, i); err != nil {
				return nil, err
			}

			results = append(results, virtualmachineimagetemplate.ImageTemplateFileCustomizer{
				Name:           pointer.To(customizer.Name),
				SourceUri:      pointer.To(customizer.FileSourceUri),
				Sha256Checksum: pointer.To(customizer.FileSha256Checksum),
				Destination:    pointer.To(customizer.FileDestinationPath),
			})

			i++
		case "PowerShell":
			if err := validateImageTemplateCustomizerInputForPowerShellType(d, customizer, i); err != nil {
				return nil, err
			}

			results = append(results, virtualmachineimagetemplate.ImageTemplatePowerShellCustomizer{
				Name:           pointer.To(customizer.Name),
				ScriptUri:      pointer.To(customizer.PowershellScriptUri),
				Sha256Checksum: pointer.To(customizer.PowershellSha256Checksum),
				Inline:         pointer.To(customizer.PowershellCommands),
				RunAsSystem:    pointer.To(customizer.PowershellRunAsSystem),
				RunElevated:    pointer.To(customizer.PowershellRunElevated),
				ValidExitCodes: pointer.To(customizer.PowershellValidExitCodes),
			})

			i++
		case "Shell":
			if err := validateImageTemplateCustomizerInputForShellType(d, customizer, i); err != nil {
				return nil, err
			}

			results = append(results, virtualmachineimagetemplate.ImageTemplateShellCustomizer{
				Name:           pointer.To(customizer.Name),
				ScriptUri:      pointer.To(customizer.ShellScriptUri),
				Sha256Checksum: pointer.To(customizer.ShellSha256Checksum),
				Inline:         pointer.To(customizer.ShellCommands),
			})

			i++
		case "WindowsRestart":
			if err := validateImageTemplateCustomizerInputForWindowsRestartType(d, i); err != nil {
				return nil, err
			}

			results = append(results, virtualmachineimagetemplate.ImageTemplateRestartCustomizer{
				Name:                pointer.To(customizer.Name),
				RestartCommand:      pointer.To(customizer.WindowsRestartCommand),
				RestartCheckCommand: pointer.To(customizer.WindowsRestartCheckCommand),
				RestartTimeout:      pointer.To(customizer.WindowsRestartTimeout),
			})

			i++
		case "WindowsUpdate":
			if err := validateImageTemplateCustomizerInputForWindowsUpdateType(d, i); err != nil {
				return nil, err
			}

			results = append(results, virtualmachineimagetemplate.ImageTemplateWindowsUpdateCustomizer{
				Name:           pointer.To(customizer.Name),
				SearchCriteria: pointer.To(customizer.WindowsUpdateSearchCriteria),
				Filters:        pointer.To(customizer.WindowsUpdateFilters),
				UpdateLimit:    pointer.To(customizer.WindowsUpdateLimit),
			})

			i++
		}
	}

	return &results, nil
}
