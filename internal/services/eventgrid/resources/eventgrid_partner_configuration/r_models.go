// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventgrid_partner_configuration

type EventGridPartnerConfigurationResourceModel struct {
	ResourceGroup                      string                 `tfschema:"resource_group_name"`
	DefaultMaximumExpirationTimeInDays int64                  `tfschema:"default_maximum_expiration_time_in_days"`
	PartnerAuthorizations              []PartnerAuthorization `tfschema:"partner_authorization"`
	Tags                               map[string]string      `tfschema:"tags"`
}

func (Resource) ModelObject() any {
	return &EventGridPartnerConfigurationResourceModel{}
}
