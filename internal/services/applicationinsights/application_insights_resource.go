// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package applicationinsights

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2019-06-01/smartdetectoralertrules"
	billing "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2015-05-01/componentfeaturesandpricingapis"
	components "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceApplicationInsights() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceApplicationInsightsCreateUpdate,
		Read:   resourceApplicationInsightsRead,
		Update: resourceApplicationInsightsCreateUpdate,
		Delete: resourceApplicationInsightsDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := components.ParseComponentID(id)
			return err
		}),

		SchemaVersion: 2,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ComponentUpgradeV0ToV1{},
			1: migration.ComponentUpgradeV1ToV2{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"application_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"web",
					"other",
					"java",
					"MobileCenter",
					"phone",
					"store",
					"ios",
					"Node.JS",
				}, false),
			},

			// NOTE: O+C A Log Analytics Workspace will be attached to the Application Insight by default, which should be computed=true
			"workspace_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: workspaces.ValidateWorkspaceID,
			},

			"retention_in_days": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				Default:  90,
				ValidateFunc: validation.IntInSlice([]int{
					30,
					60,
					90,
					120,
					180,
					270,
					365,
					550,
					730,
				}),
			},

			"sampling_percentage": {
				Type:         pluginsdk.TypeFloat,
				Optional:     true,
				Default:      100,
				ValidateFunc: validation.FloatBetween(0, 100),
			},

			"disable_ip_masking": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"tags": commonschema.Tags(),

			"daily_data_cap_in_gb": {
				Type:         pluginsdk.TypeFloat,
				Optional:     true,
				Default:      100,
				ValidateFunc: validation.FloatAtLeast(0),
			},

			"daily_data_cap_notifications_disabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"app_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"instrumentation_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"local_authentication_disabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"internet_ingestion_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"internet_query_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},
			"force_customer_storage_for_profiler": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}

	return resource
}

func resourceApplicationInsightsCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppInsights.ComponentsClient
	ruleClient := meta.(*clients.Client).Monitor.SmartDetectorAlertRulesClient
	billingClient := meta.(*clients.Client).AppInsights.BillingClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := components.NewComponentID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.ComponentsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_application_insights", id.ID())
		}
	}

	internetIngestionEnabled := components.PublicNetworkAccessTypeDisabled
	if d.Get("internet_ingestion_enabled").(bool) {
		internetIngestionEnabled = components.PublicNetworkAccessTypeEnabled
	}

	internetQueryEnabled := components.PublicNetworkAccessTypeDisabled
	if d.Get("internet_query_enabled").(bool) {
		internetQueryEnabled = components.PublicNetworkAccessTypeEnabled
	}

	applicationInsightsComponentProperties := components.ApplicationInsightsComponentProperties{
		ApplicationId:                   pointer.To(id.ComponentName),
		ApplicationType:                 components.ApplicationType(d.Get("application_type").(string)),
		SamplingPercentage:              pointer.To(d.Get("sampling_percentage").(float64)),
		DisableIPMasking:                pointer.To(d.Get("disable_ip_masking").(bool)),
		DisableLocalAuth:                pointer.To(d.Get("local_authentication_disabled").(bool)),
		PublicNetworkAccessForIngestion: pointer.To(internetIngestionEnabled),
		PublicNetworkAccessForQuery:     pointer.To(internetQueryEnabled),
		ForceCustomerStorageForProfiler: pointer.To(d.Get("force_customer_storage_for_profiler").(bool)),
	}

	if !d.IsNewResource() {
		oldWorkspaceId, newWorkspaceId := d.GetChange("workspace_id")
		if oldWorkspaceId.(string) != "" && newWorkspaceId.(string) == "" {
			return fmt.Errorf("`workspace_id` cannot be removed after set. If `workspace_id` is not specified but you encounter a diff, this might indicate a Microsoft initiated automatic migration from classic resources to workspace-based resources. If this is the case, please update `workspace_id` in your config file to the new value.")
		}
	}

	if workspaceRaw, hasWorkspaceId := d.GetOk("workspace_id"); hasWorkspaceId {
		workspaceID, err := workspaces.ParseWorkspaceID(workspaceRaw.(string))
		if err != nil {
			return err
		}
		applicationInsightsComponentProperties.WorkspaceResourceId = pointer.To(workspaceID.ID())
	}

	if v, ok := d.GetOk("retention_in_days"); ok {
		applicationInsightsComponentProperties.RetentionInDays = pointer.To(int64(v.(int)))
	}

	insightProperties := components.ApplicationInsightsComponent{
		Name:       pointer.To(id.ComponentName),
		Location:   location.Normalize(d.Get("location").(string)),
		Kind:       d.Get("application_type").(string),
		Properties: &applicationInsightsComponentProperties,
		Tags:       tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	_, err := client.ComponentsCreateOrUpdate(ctx, id, insightProperties)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	read, err := client.ComponentsGet(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	if read.Model == nil && read.Model.Id == nil {
		return fmt.Errorf("cannot read %s", id)
	}

	billingId, err := billing.ParseComponentID(id.ID())
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
		billingRead.Model.DataVolumeCap = &billing.ApplicationInsightsComponentDataVolumeCap{}
	}

	applicationInsightsComponentBillingFeatures := billing.ApplicationInsightsComponentBillingFeatures{
		CurrentBillingFeatures: billingRead.Model.CurrentBillingFeatures,
		DataVolumeCap:          billingRead.Model.DataVolumeCap,
	}

	if v, ok := d.GetOk("daily_data_cap_in_gb"); ok {
		applicationInsightsComponentBillingFeatures.DataVolumeCap.Cap = utils.Float(v.(float64))
	}

	if v, ok := d.GetOk("daily_data_cap_notifications_disabled"); ok {
		applicationInsightsComponentBillingFeatures.DataVolumeCap.StopSendNotificationWhenHitCap = utils.Bool(v.(bool))
	}

	if _, err = billingClient.ComponentCurrentBillingFeaturesUpdate(ctx, *billingId, applicationInsightsComponentBillingFeatures); err != nil {
		return fmt.Errorf("update Billing Feature for %s: %+v", id, err)
	}

	// https://github.com/hashicorp/terraform-provider-azurerm/issues/10563
	// Azure creates a rule and action group when creating this resource that are very noisy
	// We would like to delete them but deleting them just causes them to be recreated after a few minutes.
	// Instead, we'll opt to disable them here
	if d.IsNewResource() && meta.(*clients.Client).Features.ApplicationInsights.DisableGeneratedRule {
		// TODO: replace this with a StateWait func
		err = pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutCreate), func() *pluginsdk.RetryError {
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
		})
		if err != nil {
			return err
		}
	}

	d.SetId(id.ID())

	return resourceApplicationInsightsRead(d, meta)
}

func resourceApplicationInsightsRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppInsights.ComponentsClient
	billingClient := meta.(*clients.Client).AppInsights.BillingClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := components.ParseComponentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.ComponentsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	billingId, err := billing.ParseComponentID(id.ID())
	if err != nil {
		return err
	}
	billingResp, err := billingClient.ComponentCurrentBillingFeaturesGet(ctx, *billingId)
	if err != nil {
		return fmt.Errorf("retrieving Billing Features for %s: %+v", id, err)
	}

	d.Set("name", id.ComponentName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("flattening `tags`: %+v", err)
		}

		if props := model.Properties; props != nil {
			vals := map[string]string{
				"web":   "web",
				"other": "other",
			}

			if v, ok := vals[strings.ToLower(string(props.ApplicationType))]; ok {
				d.Set("application_type", v)
			} else {
				d.Set("application_type", string(props.ApplicationType))
			}
			d.Set("app_id", props.AppId)
			d.Set("instrumentation_key", props.InstrumentationKey)
			d.Set("sampling_percentage", props.SamplingPercentage)
			d.Set("disable_ip_masking", props.DisableIPMasking)
			d.Set("connection_string", props.ConnectionString)
			d.Set("local_authentication_disabled", props.DisableLocalAuth)
			d.Set("internet_ingestion_enabled", pointer.From(props.PublicNetworkAccessForIngestion) == components.PublicNetworkAccessTypeEnabled)
			d.Set("internet_query_enabled", pointer.From(props.PublicNetworkAccessForQuery) == components.PublicNetworkAccessTypeEnabled)
			d.Set("force_customer_storage_for_profiler", props.ForceCustomerStorageForProfiler)
			d.Set("retention_in_days", pointer.From(props.RetentionInDays))
			workspaceId := ""
			if v := props.WorkspaceResourceId; v != nil {
				id, err := workspaces.ParseWorkspaceIDInsensitively(*v)
				if err != nil {
					return err
				}
				workspaceId = id.ID()
			}
			d.Set("workspace_id", workspaceId)
		}
	}

	if model := billingResp.Model; model != nil {
		if props := model.DataVolumeCap; props != nil {
			d.Set("daily_data_cap_in_gb", props.Cap)
			d.Set("daily_data_cap_notifications_disabled", props.StopSendNotificationWhenHitCap)
		}
	}

	return nil
}

func resourceApplicationInsightsDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppInsights.ComponentsClient
	ruleClient := meta.(*clients.Client).Monitor.SmartDetectorAlertRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := components.ParseComponentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.ComponentsDelete(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	// if disable_generated_rule=true, the generated rule is not automatically deleted.
	if meta.(*clients.Client).Features.ApplicationInsights.DisableGeneratedRule {
		ruleName := fmt.Sprintf("Failure Anomalies - %s", id.ComponentName)
		ruleId := smartdetectoralertrules.NewSmartDetectorAlertRuleID(id.SubscriptionId, id.ResourceGroupName, ruleName)
		deleteResp, deleteErr := ruleClient.Delete(ctx, ruleId)
		if deleteErr != nil && !response.WasNotFound(deleteResp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", ruleId, deleteErr)
		}
	}

	return err
}
