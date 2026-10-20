// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_trigger_http_request

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflowtriggers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogicAppTriggerHttpRequestCreateUpdate,
		Read:   resourceLogicAppTriggerHttpRequestRead,
		Update: resourceLogicAppTriggerHttpRequestCreateUpdate,
		Delete: resourceLogicAppTriggerHttpRequestDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := workflowtriggers.ParseTriggerID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v any) error {
			relativePath := diff.Get("relative_path").(string)
			if relativePath != "" {
				method := diff.Get("method").(string)
				if method == "" {
					return fmt.Errorf("`method` must be specified when `relative_path` is set")
				}
			}

			return nil
		}),

		Schema: logicAppTriggerHttpRequestSchema(),
	}
}
