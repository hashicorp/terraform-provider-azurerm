// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package oracle

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/cloudvmclusters"
)

func FlattenFileSystemConfigurationDetails(input *[]cloudvmclusters.FileSystemConfigurationDetails) []FileSystemConfigurationModel {
	output := make([]FileSystemConfigurationModel, 0)
	if input != nil {
		for _, item := range *input {
			output = append(output, FileSystemConfigurationModel{
				MountPoint: pointer.From(item.MountPoint),
				SizeInGb:   pointer.From(item.FileSystemSizeGb),
			})
		}
	}
	return output
}

func ExpandFileSystemConfiguration(fileSystemConfigurations []FileSystemConfigurationModel) *[]cloudvmclusters.FileSystemConfigurationDetails {
	properties := make([]cloudvmclusters.FileSystemConfigurationDetails, 0)
	for _, item := range fileSystemConfigurations {
		// We need to skip mount points not allowed to resize
		if item.MountPoint != "reserved" && item.MountPoint != "swap" && item.MountPoint != "/var/log/audit" {
			properties = append(properties, cloudvmclusters.FileSystemConfigurationDetails{
				MountPoint:       pointer.To(item.MountPoint),
				FileSystemSizeGb: pointer.To(item.SizeInGb),
			})
		}
	}
	return &properties
}
