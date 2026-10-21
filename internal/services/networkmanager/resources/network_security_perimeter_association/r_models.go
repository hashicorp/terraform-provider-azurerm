// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_security_perimeter_association

type NetworkSecurityPerimeterAssociationResourceModel struct {
	Name       string `tfschema:"name"`
	ProfileId  string `tfschema:"network_security_perimeter_profile_id"`
	ResourceId string `tfschema:"resource_id"`
	AccessMode string `tfschema:"access_mode"`
}

func (Resource) ModelObject() any {
	return &NetworkSecurityPerimeterAssociationResourceModel{}
}
