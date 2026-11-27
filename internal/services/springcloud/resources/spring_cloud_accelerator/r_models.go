// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_accelerator

type SpringCloudAcceleratorModel struct {
	Name                 string `tfschema:"name"`
	SpringCloudServiceId string `tfschema:"spring_cloud_service_id"`
}

func (s Resource) ModelObject() any {
	return &SpringCloudAcceleratorModel{}
}
