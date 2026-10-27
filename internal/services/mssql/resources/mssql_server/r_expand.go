// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_server

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/serverazureadadministrators"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/servers"
)

func expandMsSqlServerAADOnlyAuthentication(input []any) bool {
	if len(input) == 0 || input[0] == nil {
		return false
	}

	admin := input[0].(map[string]any)

	if v, ok := admin["azuread_authentication_only"]; ok && v != nil {
		return v.(bool)
	}

	return false
}

func expandMsSqlServerAdministrator(input []any) *serverazureadadministrators.ServerAzureADAdministrator {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]any)

	adminProps := serverazureadadministrators.ServerAzureADAdministrator{
		Properties: &serverazureadadministrators.AdministratorProperties{
			AdministratorType: pointer.To(serverazureadadministrators.AdministratorTypeActiveDirectory),
			Login:             v["login_username"].(string),
			Sid:               v["object_id"].(string),
		},
	}

	if t, ok := v["tenant_id"]; ok && t != "" {
		adminProps.Properties.TenantId = pointer.To(t.(string))
	}

	return pointer.To(adminProps)
}

func expandMsSqlServerAdministrators(input []any) *servers.ServerExternalAdministrator {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	admin := input[0].(map[string]any)
	sid := admin["object_id"].(string)

	adminParams := servers.ServerExternalAdministrator{
		AdministratorType: pointer.To(servers.AdministratorTypeActiveDirectory),
		Login:             pointer.To(admin["login_username"].(string)),
		Sid:               pointer.To(sid),
	}

	if v, ok := admin["tenant_id"]; ok && v != "" {
		adminParams.TenantId = pointer.To(v.(string))
	}

	if v, ok := admin["azuread_authentication_only"]; ok && v != "" {
		adminParams.AzureADOnlyAuthentication = pointer.To(v.(bool))
	}

	return &adminParams
}
