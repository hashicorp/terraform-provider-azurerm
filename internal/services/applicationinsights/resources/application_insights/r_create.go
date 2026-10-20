// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2019-06-01/smartdetectoralertrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2015-05-01/componentfeaturesandpricingapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApplicationInsightsCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppInsights.ComponentsClient
	ruleClient := meta.(*clients.Client).Monitor.SmartDetectorAlertRulesClient
	billingClient := meta.(*clients.Client).AppInsights.BillingClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := componentsapis.NewComponentID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.ComponentsGet(ctx, id)
		if !response.WasNotFound(existing.HttpResponse) {
			if err != nil {
				return fmt.Errorf("checking for presence of existing %s: %v", id, err)
			}
			return tf.ImportAsExistsError("azurerm_application_insights", id.ID())
		}
	}

	internetIngestionEnabled := componentsapis.PublicNetworkAccessTypeDisabled
	if d.Get("internet_ingestion_enabled").(bool) {
		internetIngestionEnabled = componentsapis.PublicNetworkAccessTypeEnabled
	}

	internetQueryEnabled := componentsapis.PublicNetworkAccessTypeDisabled
	if d.Get("internet_query_enabled").(bool) {
		internetQueryEnabled = componentsapis.PublicNetworkAccessTypeEnabled
	}

	applicationInsightsComponentProperties := componentsapis.ApplicationInsightsComponentProperties{
		ApplicationId:                   pointer.To(id.ComponentName),
		ApplicationType:                 componentsapis.ApplicationType(d.Get("application_type").(string)),
		SamplingPercentage:              pointer.To(d.Get("sampling_percentage").(float64)),
		DisableIPMasking:                pointer.To(!d.Get("ip_masking_enabled").(bool)),
		DisableLocalAuth:                pointer.To(!d.Get("local_authentication_enabled").(bool)),
		PublicNetworkAccessForIngestion: pointer.To(internetIngestionEnabled),
		PublicNetworkAccessForQuery:     pointer.To(internetQueryEnabled),
		ForceCustomerStorageForProfiler: pointer.To(d.Get("force_customer_storage_for_profiler").(bool)),
	}

	if workspaceRaw, ok := d.GetOk("workspace_id"); ok {
		workspaceID, err := workspaces.ParseWorkspaceID(workspaceRaw.(string))
		if err != nil {
			return err
		}
		applicationInsightsComponentProperties.WorkspaceResourceId = pointer.To(workspaceID.ID())
	}

	if v, ok := d.GetOk("retention_in_days"); ok {
		applicationInsightsComponentProperties.RetentionInDays = pointer.To(int64(v.(int)))
	}

	insightProperties := componentsapis.ApplicationInsightsComponent{
		Name:       pointer.To(id.ComponentName),
		Location:   location.Normalize(d.Get("location").(string)),
		Kind:       d.Get("application_type").(string),
		Properties: &applicationInsightsComponentProperties,
		Tags:       tags.Expand(d.Get("tags").(map[string]any)),
	}

	if _, err := client.ComponentsCreateOrUpdate(ctx, id, insightProperties); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	read, err := client.ComponentsGet(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	if read.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	if read.Model.Id == nil {
		return fmt.Errorf("retrieving %s: `id` was nil", id)
	}

	billingId, err := componentfeaturesandpricingapis.ParseComponentID(id.ID())
	if err != nil {
		return err
	}
	billingRead, err := billingClient.ComponentCurrentBillingFeaturesGet(ctx, *billingId)
	if err != nil {
		return fmt.Errorf("retrieving Billing Features for %s: %+v", id, err)
	}

	if billingRead.Model == nil {
		return fmt.Errorf("model is nil for billing features")
	}

	if billingRead.Model.DataVolumeCap == nil {
		billingRead.Model.DataVolumeCap = &componentfeaturesandpricingapis.ApplicationInsightsComponentDataVolumeCap{}
	}

	applicationInsightsComponentBillingFeatures := componentfeaturesandpricingapis.ApplicationInsightsComponentBillingFeatures{
		CurrentBillingFeatures: billingRead.Model.CurrentBillingFeatures,
		DataVolumeCap:          billingRead.Model.DataVolumeCap,
	}

	if v, ok := d.GetOk("daily_data_cap_in_gb"); ok {
		applicationInsightsComponentBillingFeatures.DataVolumeCap.Cap = pointer.To(v.(float64))
	}

	applicationInsightsComponentBillingFeatures.DataVolumeCap.StopSendNotificationWhenHitCap = pointer.To(!d.Get("daily_data_cap_notifications_enabled").(bool))

	if _, err = billingClient.ComponentCurrentBillingFeaturesUpdate(ctx, *billingId, applicationInsightsComponentBillingFeatures); err != nil {
		return fmt.Errorf("update Billing Feature for %s: %+v", id, err)
	}

	// https://github.com/hashicorp/terraform-provider-azurerm/issues/10563
	// Azure creates a rule and action group when creating this resource that are very noisy
	// We would like to delete them but deleting them just causes them to be recreated after a few minutes.
	// Instead, we'll opt to disable them here
	if meta.(*clients.Client).Features.ApplicationInsights.DisableGeneratedRule {
		// TODO: replace this with a StateWait func
		if err = pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutCreate), func() *pluginsdk.RetryError {
			time.Sleep(30 * time.Second)
			ruleName := fmt.Sprintf("Failure Anomalies - %s", id.ComponentName)
			ruleId := smartdetectoralertrules.NewSmartDetectorAlertRuleID(id.SubscriptionId, id.ResourceGroupName, ruleName)
			result, err := ruleClient.Get(ctx, ruleId, smartdetectoralertrules.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(result.HttpResponse) {
					return pluginsdk.RetryableError(fmt.Errorf("expected %s to be created but was not found, retrying", ruleId))
				}
				return pluginsdk.NonRetryableError(fmt.Errorf("making Read request for %s: %+v", ruleId, err))
			}

			if model := result.Model; model != nil {
				if props := model.Properties; props != nil {
					props.State = smartdetectoralertrules.AlertRuleStateDisabled
					updateRuleResult, err := ruleClient.CreateOrUpdate(ctx, ruleId, *model)
					if err != nil {
						if !response.WasNotFound(updateRuleResult.HttpResponse) {
							return pluginsdk.NonRetryableError(fmt.Errorf("issuing disable request for %s: %+v", ruleId, err))
						}
					}
				}
			}

			return nil
		}); err != nil {
			return err
		}
	}

	return resourceApplicationInsightsRead(d, meta)
}
