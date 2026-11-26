// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_center_dev_box_definition

func (r Resource) ModelObject() any {
	return &DevCenterDevBoxDefinitionResourceModel{}
}

type DevCenterDevBoxDefinitionResourceModel struct {
	Name             string            `tfschema:"name"`
	Location         string            `tfschema:"location"`
	DevCenterId      string            `tfschema:"dev_center_id"`
	HibernateSupport bool              `tfschema:"hibernate_support_enabled"`
	ImageReferenceId string            `tfschema:"image_reference_id"`
	SkuName          string            `tfschema:"sku_name"`
	Tags             map[string]string `tfschema:"tags"`
}
