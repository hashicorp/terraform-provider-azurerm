// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package static_web_app_custom_domain

type StaticWebAppCustomDomainResourceModel struct {
	DomainName      string `tfschema:"domain_name"`
	StaticSiteId    string `tfschema:"static_web_app_id"`
	ValidationType  string `tfschema:"validation_type"`
	ValidationToken string `tfschema:"validation_token"`
}

func (r Resource) ModelObject() any {
	return &Resource{}
}
