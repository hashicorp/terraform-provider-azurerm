// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package gallery_application_version

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryapplicationversions"
)

func flattenGalleryApplicationVersionManageAction(input *galleryapplicationversions.UserArtifactManage) []ManageAction {
	if input == nil {
		return []ManageAction{}
	}

	output := make([]ManageAction, 0)

	obj := ManageAction{
		Install: input.Install,
		Remove:  input.Remove,
	}
	if input.Update != nil {
		obj.Update = *input.Update
	}
	output = append(output, obj)

	return output
}

func flattenGalleryApplicationVersionSource(input galleryapplicationversions.UserArtifactSource) []Source {
	out := Source{
		MediaLink: input.MediaLink,
	}
	if input.DefaultConfigurationLink != nil {
		out.DefaultConfigurationLink = *input.DefaultConfigurationLink
	}
	return []Source{
		out,
	}
}

func flattenGalleryApplicationVersionTargetRegion(input *[]galleryapplicationversions.TargetRegion) []TargetRegion {
	results := make([]TargetRegion, 0)

	for _, item := range *input {
		obj := TargetRegion{
			Name: location.Normalize(item.Name),
		}

		if item.ExcludeFromLatest != nil {
			obj.ExcludeFromLatest = *item.ExcludeFromLatest
		}

		if item.RegionalReplicaCount != nil {
			obj.RegionalReplicaCount = *item.RegionalReplicaCount
		}

		if item.StorageAccountType != nil {
			obj.StorageAccountType = string(*item.StorageAccountType)
		}

		results = append(results, obj)
	}

	return results
}
