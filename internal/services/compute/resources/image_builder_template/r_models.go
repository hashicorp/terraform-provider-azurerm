// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package image_builder_template

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

func (Resource) ModelObject() any {
	return &ImageBuilderTemplateResourceModel{}
}

type ImageBuilderTemplateResourceModel struct {
	Name                       string                                    `tfschema:"name"`
	ResourceGroupName          string                                    `tfschema:"resource_group_name"`
	Location                   string                                    `tfschema:"location"`
	Identity                   []identity.ModelUserAssigned              `tfschema:"identity"`
	BuildTimeoutMinutes        int64                                     `tfschema:"build_timeout_minutes"`
	Customizer                 []ImageBuilderTemplateCustomizer          `tfschema:"customizer"`
	DiskSizeGb                 int64                                     `tfschema:"disk_size_gb"`
	Distributions              []ImageBuilderTemplateDistributions       `tfschema:"distributions"`
	Size                       string                                    `tfschema:"size"`
	SourceManagedImageId       string                                    `tfschema:"source_managed_image_id"`
	SourcePlatformImage        []ImageBuilderTemplateSourcePlatformImage `tfschema:"source_platform_image"`
	SourceSharedImageVersionId string                                    `tfschema:"source_shared_image_version_id"`
	SubnetId                   string                                    `tfschema:"subnet_id"`
	Tags                       map[string]string                         `tfschema:"tags"`
}
