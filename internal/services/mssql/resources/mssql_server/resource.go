// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_server

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -test-expect-non-empty

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Manages a Microsoft SQL Azure Database Server.",
		Create:      resourceMsSqlServerCreate,
		Read:        resourceMsSqlServerRead,
		Update:      resourceMsSqlServerUpdate,
		Delete:      resourceMsSqlServerDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&commonids.SqlServerId{}),
		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&commonids.SqlServerId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: mssqlServerSchema(),

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.CustomizeDiffShim(msSqlMinimumTLSVersionDiff),

			pluginsdk.CustomizeDiffShim(msSqlPasswordChangeWhenAADAuthOnly),
		),
	}
}

func msSqlMinimumTLSVersionDiff(ctx context.Context, d *pluginsdk.ResourceDiff, _ any) (err error) {
	old, new := d.GetChange("minimum_tls_version")
	// todo remove `old != "None"` when https://github.com/Azure/azure-rest-api-specs/issues/24348 is addressed
	if old != "" && old != "None" && old != "Disabled" && new == "Disabled" {
		err = fmt.Errorf("`minimum_tls_version` cannot be removed once set, please set a valid value for this property")
	}
	return
}

func msSqlPasswordChangeWhenAADAuthOnly(ctx context.Context, d *pluginsdk.ResourceDiff, _ any) error {
	if d.Id() == "" {
		return nil
	}

	_, newVal := d.GetChange("azuread_administrator.0.azuread_authentication_only")
	if !newVal.(bool) {
		return nil
	}

	if d.HasChange("administrator_login_password") {
		_, newPwd := d.GetChange("administrator_login_password")
		if !d.NewValueKnown("administrator_login_password") || newPwd.(string) != "" {
			return fmt.Errorf("`administrator_login_password` cannot be changed when `azuread_administrator.0.azuread_authentication_only = true`")
		}
	}

	if d.HasChange("administrator_login_password_wo_version") {
		_, newVersion := d.GetChange("administrator_login_password_wo_version")
		if !d.NewValueKnown("administrator_login_password_wo_version") || newVersion.(int) != 0 {
			return fmt.Errorf("`administrator_login_password_wo` cannot be changed when `azuread_administrator.0.azuread_authentication_only = true`")
		}
	}

	return nil
}
