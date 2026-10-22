// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_elasticpool

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMsSqlElasticPoolCreateUpdate,
		Read:   resourceMsSqlElasticPoolRead,
		Update: resourceMsSqlElasticPoolCreateUpdate,
		Delete: resourceMsSqlElasticPoolDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&commonids.SqlElasticPoolId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&commonids.SqlElasticPoolId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: mssqlElasticpoolSchema(),

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			func(ctx context.Context, diff *pluginsdk.ResourceDiff, v any) error {
				if err := helper.MSSQLElasticPoolValidateSKU(diff); err != nil {
					return err
				}

				if v, ok := diff.GetOk("high_availability_replica_count"); ok && v.(int) > 0 {
					skuRaw := diff.Get("sku").([]any)
					if len(skuRaw) > 0 {
						sku := skuRaw[0].(map[string]any)
						tier := sku["tier"].(string)
						if !strings.EqualFold(tier, "Hyperscale") {
							return fmt.Errorf("`high_availability_replica_count` can only be set when `sku.tier` is `Hyperscale`, got %q", tier)
						}
					}
				}

				return nil
			},

			pluginsdk.ForceNewIfChange("enclave_type", func(ctx context.Context, old, new, _ any) bool {
				// enclave_type cannot be removed once it has been set
				// but can be changed between VBS and Default
				if old.(string) != "" && new.(string) == "" {
					return true
				}

				return false
			}),
		),
	}
}
