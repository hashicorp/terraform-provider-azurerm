// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_partner_registration

type EventGridPartnerRegistrationResourceModel struct {
	Name                  string            `tfschema:"name"`
	ResourceGroup         string            `tfschema:"resource_group_name"`
	PartnerRegistrationID string            `tfschema:"partner_registration_id"`
	Tags                  map[string]string `tfschema:"tags"`
}

func (Resource) ModelObject() any {
	return &EventGridPartnerRegistrationResourceModel{}
}
