// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_security_perimeter_profile

type NetworkSecurityPerimeterProfileResourceModel struct {
	Name        string `tfschema:"name"`
	PerimeterId string `tfschema:"network_security_perimeter_id"`
}

func (Resource) ModelObject() any {
	return &NetworkSecurityPerimeterProfileResourceModel{}
}
