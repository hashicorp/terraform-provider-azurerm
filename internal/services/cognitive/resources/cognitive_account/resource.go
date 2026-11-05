// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_account

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/cognitiveservicesaccounts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

const azureCognitiveAccountResourceName = "azurerm_cognitive_account"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCognitiveAccountCreate,
		Read:   resourceCognitiveAccountRead,
		Update: resourceCognitiveAccountUpdate,
		Delete: resourceCognitiveAccountDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingIdentity(&cognitiveservicesaccounts.AccountId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&cognitiveservicesaccounts.AccountId{}),
		},

		Schema: cognitiveAccountSchema(),

		CustomizeDiff: func(ctx context.Context, d *pluginsdk.ResourceDiff, i any) error {
			kind := d.Get("kind").(string)

			if d.Get("project_management_enabled").(bool) {
				if kind != "AIServices" {
					return errors.New("`project_management_enabled` can only be set to `true` when `kind` is set to `AIServices`")
				}

				if len(d.Get("identity").([]any)) == 0 {
					return errors.New("for `project_management_enabled` to be set to `true`, a managed identity must be assigned. Please configure the `identity` block")
				}

				if d.HasChange("customer_managed_key") && len(d.Get("customer_managed_key").([]any)) == 0 {
					if err := d.ForceNew("customer_managed_key"); err != nil {
						return err
					}
				}
			} else if d.HasChange("project_management_enabled") && kind != "OpenAI" {
				if err := d.ForceNew("project_management_enabled"); err != nil {
					return err
				}
			}

			if d.Get("dynamic_throttling_enabled").(bool) && slices.Contains([]string{"OpenAI", "AIServices"}, kind) {
				return errors.New("`dynamic_throttling_enabled` is currently not supported when `kind` is set to `OpenAI` or `AIServices`")
			}

			if bypass, ok := d.GetOk("network_acls.0.bypass"); ok && bypass != "" && !slices.Contains([]string{"OpenAI", "AIServices", "TextAnalytics"}, kind) {
				return fmt.Errorf("`network_acls.bypass` cannot be set when `kind` is set to `%s`", kind)
			}

			networkInjection := d.Get("network_injection").([]any)
			if len(networkInjection) > 0 && networkInjection[0] != nil {
				if kind != "AIServices" {
					return errors.New("the `network_injection` block is only supported when `kind` is set to `AIServices`")
				}
			}

			if d.HasChange("custom_subdomain_name") {
				old, _ := d.GetChange("custom_subdomain_name")
				if old != nil && old != "" {
					if err := d.ForceNew("custom_subdomain_name"); err != nil {
						return err
					}
				}
			}

			if d.HasChange("kind") {
				old, new := d.GetChange("kind")

				// Only allow changing `kind` from/to `OpenAI` or `AIServices`, force new for all others
				if !slices.Contains([]string{"OpenAI", "AIServices"}, new.(string)) || !slices.Contains([]string{"OpenAI", "AIServices"}, old.(string)) {
					if err := d.ForceNew("kind"); err != nil {
						return err
					}
				}
			}

			return nil
		},
	}
}
