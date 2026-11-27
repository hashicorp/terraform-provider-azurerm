// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dynatrace_tag_rules

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dynatrace/2023-04-27/monitors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dynatrace/2023-04-27/tagrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	dynatraceClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/dynatrace/client"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := dynatraceClient.From(metadata.Client).TagRulesClient
			subscriptionId := account.From(metadata.Client).SubscriptionId

			var model TagRulesResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			monitorId, err := monitors.ParseMonitorID(model.Monitor)
			id := tagrules.NewTagRuleID(subscriptionId, monitorId.ResourceGroupName, monitorId.MonitorName, model.Name)
			if err != nil {
				return err
			}

			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			tagRulesProps := tagrules.MonitoringTagRulesProperties{
				LogRules:    ExpandLogRule(model.LogRules),
				MetricRules: ExpandMetricRules(model.MetricRules),
			}
			tagRules := tagrules.TagRule{
				Name:       &model.Name,
				Properties: tagRulesProps,
			}

			if _, err := client.CreateOrUpdate(ctx, id, tagRules); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}
