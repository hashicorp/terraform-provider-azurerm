// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_diagnostic

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/diagnostic"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementDiagnosticCreateUpdate,
		Read:   resourceApiManagementDiagnosticRead,
		Update: resourceApiManagementDiagnosticCreateUpdate,
		Delete: resourceApiManagementDiagnosticDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := diagnostic.ParseDiagnosticID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: apiManagementDiagnosticSchema(),

		CustomizeDiff: func(ctx context.Context, d *pluginsdk.ResourceDiff, i any) error {
			if _, n := d.GetChange("operation_name_format"); n != "" {
				if d.Get("identifier") != "applicationinsights" {
					return fmt.Errorf("`operation_name_format` cannot be set when `identifier` is not `applicationinsights`")
				}
			}

			return nil
		},
	}
}
