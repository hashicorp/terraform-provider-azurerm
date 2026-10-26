// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package private_endpoint_application_security_group_association

type PrivateEndpointApplicationSecurityGroupAssociationModel struct {
	PrivateEndpointId          string `tfschema:"private_endpoint_id"`
	ApplicationSecurityGroupId string `tfschema:"application_security_group_id"`
}

func (p Resource) ModelObject() any {
	return &PrivateEndpointApplicationSecurityGroupAssociationModel{}
}
