// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package datadog_monitor

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2025-06-11/datadogmonitorresources"
)

func expandMonitorIdentityProperties(input []any) *datadogmonitorresources.IdentityProperties {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]any)
	return &datadogmonitorresources.IdentityProperties{
		// @tombuildsstuff: this should be normalized in Pandora to a common identity type? is this a Swagger bag with SA & UA omitted?
		Type: pointer.ToEnum[datadogmonitorresources.ManagedIdentityTypes](v["type"].(string)),
	}
}

func expandMonitorOrganizationProperties(input []any) *datadogmonitorresources.DatadogOrganizationProperties {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]any)
	return &datadogmonitorresources.DatadogOrganizationProperties{
		LinkingAuthCode: pointer.To(v["linking_auth_code"].(string)),
		LinkingClientId: pointer.To(v["linking_client_id"].(string)),
		RedirectUri:     pointer.To(v["redirect_uri"].(string)),
		ApiKey:          pointer.To(v["api_key"].(string)),
		ApplicationKey:  pointer.To(v["application_key"].(string)),
		EnterpriseAppId: pointer.To(v["enterprise_app_id"].(string)),
	}
}

func expandMonitorUserInfo(input []any) *datadogmonitorresources.UserInfo {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]any)
	return &datadogmonitorresources.UserInfo{
		Name:         pointer.To(v["name"].(string)),
		EmailAddress: pointer.To(v["email"].(string)),
		PhoneNumber:  pointer.To(v["phone_number"].(string)),
	}
}
