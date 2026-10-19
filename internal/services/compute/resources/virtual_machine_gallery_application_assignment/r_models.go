// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_gallery_application_assignment

type VirtualMachineGalleryApplicationAssignmentResourceResourceModel struct {
	GalleryApplicationVersionId string `tfschema:"gallery_application_version_id"`
	VirtualMachineId            string `tfschema:"virtual_machine_id"`
	ConfigurationBlobUri        string `tfschema:"configuration_blob_uri"`
	Order                       int64  `tfschema:"order"`
	Tag                         string `tfschema:"tag"`
}

func (r Resource) ModelObject() any {
	return &VirtualMachineGalleryApplicationAssignmentResourceResourceModel{}
}
