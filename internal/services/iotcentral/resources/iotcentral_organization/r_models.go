// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iotcentral_organization

type IotCentralOrganizationModel struct {
	IotCentralApplicationId string `tfschema:"iotcentral_application_id"`
	OrganizationId          string `tfschema:"organization_id"`
	DisplayName             string `tfschema:"display_name"`
	ParentOrganizationId    string `tfschema:"parent_organization_id"`
}

func (r Resource) ModelObject() any {
	return &IotCentralOrganizationModel{}
}
