// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package service_plan

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/preflight"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ServicePlanName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"sku_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice(
				helpers.AllKnownServicePlanSkus(),
				false,
			),
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"os_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(OSTypeLinux),
				string(OSTypeWindows),
				string(OSTypeWindowsContainer),
			}, false),
		},

		"app_service_environment_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateAppServiceEnvironmentID,
		},

		"per_site_scaling_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"worker_count": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.IntAtLeast(1),
		},

		"premium_plan_auto_scale_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"maximum_elastic_worker_count": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.IntAtLeast(0),
		},

		"zone_balancing_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"kind": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"reserved": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},
	}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff

			if metadata.Client.Features.EnhancedValidation.PreflightEnabled {
				// Only perform preflight validation if there are changes. This avoids validation failures and
				// additional API calls for resources that are unchanged between plan invocations
				if len(rd.GetChangedKeysPrefix("")) > 0 || rd.Id() == "" {
					var model ServicePlanModel
					if err := metadata.DecodeDiff(&model); err != nil {
						return err
					}

					req, err := expandCreateForServicePlan(model)
					if err != nil {
						return err
					}

					id := commonids.NewAppServicePlanID(metadata.Client.Account.SubscriptionId, model.ResourceGroup, model.Name)

					preflightValidate, err := preflight.NewValidationRequest(pointer.To(model.Location), pointer.To(id), "2023-12-01", req)
					if err != nil {
						return fmt.Errorf("constructing preflight validation request: %w", err)
					}

					if err = preflightValidate.ValidateResource(ctx, metadata); err != nil {
						return err
					}
				}
			}

			servicePlanSku := rd.Get("sku_name").(string)
			_, newAutoScaleEnabled := rd.GetChange("premium_plan_auto_scale_enabled")
			_, newEcValue := rd.GetChange("maximum_elastic_worker_count")
			if rd.HasChange("premium_plan_auto_scale_enabled") {
				if !helpers.PlanIsPremium(servicePlanSku) && newAutoScaleEnabled.(bool) {
					return fmt.Errorf("`premium_plan_auto_scale_enabled` can only be set for premium app service plans")
				}
			}

			if rd.HasChange("maximum_elastic_worker_count") && newEcValue.(int) > 1 {
				if !helpers.PlanSupportsScaleOut(servicePlanSku) && helpers.PlanIsPremium(servicePlanSku) && !newAutoScaleEnabled.(bool) {
					return fmt.Errorf("`maximum_elastic_worker_count` can only be specified with Elastic Premium Skus, or with Premium Skus when `premium_plan_auto_scale_enabled` is set to `true`")
				}
			}

			// Only specific SKUs support zone balancing/redundancy
			if rd.Get("zone_balancing_enabled").(bool) {
				if !helpers.PlanSupportsZoneBalancing(servicePlanSku) {
					return fmt.Errorf("`zone_balancing_enabled` cannot be set to `true` when sku tier is `%s`", servicePlanSku)
				}
			}

			o, n := rd.GetChange("zone_balancing_enabled")
			if o.(bool) != n.(bool) {
				// Changing `zone_balancing_enabled` from `false` to `true` requires the capacity of the sku to be greater than `1`.
				if !o.(bool) && n.(bool) && rd.Get("worker_count").(int) < 2 {
					if err := metadata.ResourceDiff.ForceNew("zone_balancing_enabled"); err != nil {
						return err
					}
				}
			}

			return nil
		},
	}
}
