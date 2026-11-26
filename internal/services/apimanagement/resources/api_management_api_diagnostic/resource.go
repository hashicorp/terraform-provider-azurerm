// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_api_diagnostic

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apidiagnostic"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Manages an API Management Service API Diagnostics Logs.",
		Create:      resourceApiManagementApiDiagnosticCreateUpdate,
		Read:        resourceApiManagementApiDiagnosticRead,
		Update:      resourceApiManagementApiDiagnosticCreateUpdate,
		Delete:      resourceApiManagementApiDiagnosticDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := apidiagnostic.ParseApiDiagnosticID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: apiManagementApiDiagnosticSchema(),

		CustomizeDiff: func(ctx context.Context, d *pluginsdk.ResourceDiff, i any) error {
			format := d.GetRawConfig().GetAttr("operation_name_format")
			if !format.IsNull() && d.Get("identifier") != "applicationinsights" {
				return fmt.Errorf("`operation_name_format` cannot be set when `identifier` is not `applicationinsights`")
			}

			if v := d.GetRawConfig().AsValueMap()["metrics_enabled"]; !v.IsNull() && d.Get("identifier") != "applicationinsights" {
				return fmt.Errorf("`metrics_enabled` can only be set when `identifier` is set to `applicationinsights`")
			}

			return nil
		},
	}
}
