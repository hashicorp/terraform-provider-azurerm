// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package datadog_monitor

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/monitorsresource"
)

func expandMonitorIdentityProperties(input []any) *monitorsresource.IdentityProperties {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]any)
	return &monitorsresource.IdentityProperties{
		// @tombuildsstuff: this should be normalized in Pandora to a common identity type? is this a Swagger bag with SA & UA omitted?
		Type: pointer.ToEnum[monitorsresource.ManagedIdentityTypes](v["type"].(string)),
	}
}

func expandMonitorOrganizationProperties(input []any) *monitorsresource.DatadogOrganizationProperties {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]any)
	return &monitorsresource.DatadogOrganizationProperties{
		LinkingAuthCode: pointer.To(v["linking_auth_code"].(string)),
		LinkingClientId: pointer.To(v["linking_client_id"].(string)),
		RedirectUri:     pointer.To(v["redirect_uri"].(string)),
		ApiKey:          pointer.To(v["api_key"].(string)),
		ApplicationKey:  pointer.To(v["application_key"].(string)),
		EnterpriseAppId: pointer.To(v["enterprise_app_id"].(string)),
	}
}

func expandMonitorUserInfo(input []any) *monitorsresource.UserInfo {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]any)
	return &monitorsresource.UserInfo{
		Name:         pointer.To(v["name"].(string)),
		EmailAddress: pointer.To(v["email"].(string)),
		PhoneNumber:  pointer.To(v["phone_number"].(string)),
	}
}
