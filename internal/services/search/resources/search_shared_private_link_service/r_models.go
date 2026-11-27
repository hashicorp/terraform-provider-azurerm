// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package search_shared_private_link_service

type SharedPrivateLinkServiceModel struct {
	Name             string `tfschema:"name"`
	SearchServiceId  string `tfschema:"search_service_id"`
	SubResourceName  string `tfschema:"subresource_name"`
	TargetResourceId string `tfschema:"target_resource_id"`
	RequestMessage   string `tfschema:"request_message"`
	Status           string `tfschema:"status"`
}

func (r Resource) ModelObject() any {
	return &SharedPrivateLinkServiceModel{}
}
