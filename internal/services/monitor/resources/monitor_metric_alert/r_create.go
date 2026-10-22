// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_metric_alert

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2018-03-01/metricalerts"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMonitorMetricAlertCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Monitor.MetricAlertsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := metricalerts.NewMetricAlertID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing Monitor %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_monitor_metric_alert", id.ID())
			}
		}
	}

	enabled := d.Get("enabled").(bool)
	autoMitigate := d.Get("auto_mitigate").(bool)
	description := d.Get("description").(string)
	scopesRaw := d.Get("scopes").(*pluginsdk.Set).List()
	severity := d.Get("severity").(int)
	frequency := d.Get("frequency").(string)
	windowSize := d.Get("window_size").(string)
	actionRaw := d.Get("action").(*pluginsdk.Set).List()
	targetResourceType := d.Get("target_resource_type").(string)
	targetResourceLocation := d.Get("target_resource_location").(string)

	t := d.Get("tags").(map[string]any)

	// The criteria type of "old" resource is `MetricAlertSingleResourceMultipleMetricCriteria` (rather than `MetricAlertMultipleResourceMultipleMetricCriteria`).
	// We need to keep using that type in order to keep backward compatibility. Otherwise, changing the criteria type will cause error as reported in issue:
	// https://github.com/hashicorp/terraform-provider-azurerm/issues/7910
	var isLegacy bool
	if !d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			return fmt.Errorf("retrieving Monitor %s: %+v", id, err)
		}
		if existing.Model == nil || existing.Model.Properties.Criteria == nil {
			return fmt.Errorf("unexpected nil properties of Monitor %s", id)
		}
		_, isLegacy = existing.Model.Properties.Criteria.(metricalerts.MetricAlertSingleResourceMultipleMetricCriteria)
	}

	criteria, err := expandMonitorMetricAlertCriteria(d, isLegacy)
	if err != nil {
		return fmt.Errorf(`expanding criteria: %+v`, err)
	}

	parameters := metricalerts.MetricAlertResource{
		Location: location.Normalize("Global"),
		Properties: metricalerts.MetricAlertProperties{
			Enabled:              enabled,
			AutoMitigate:         pointer.To(autoMitigate),
			Description:          pointer.To(description),
			Severity:             int64(severity),
			EvaluationFrequency:  frequency,
			WindowSize:           windowSize,
			Scopes:               helpers.ExpandStringValues(scopesRaw),
			Criteria:             criteria,
			Actions:              expandMonitorMetricAlertAction(actionRaw),
			TargetResourceType:   pointer.To(targetResourceType),
			TargetResourceRegion: pointer.To(targetResourceLocation),
		},
		Tags: pluginsdk.ExpandPtrMapStringString(t),
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating or updating Monitor %s: %+v", id, err)
	}

	// Monitor Metric Alert API would return 404 while creating multiple Monitor Metric Alerts and get each resource immediately once it's created successfully in parallel.
	// Tracked by this issue: https://github.com/Azure/azure-rest-api-specs/issues/10973
	log.Printf("[DEBUG] Waiting for %s to be created", id)
	poller := custompollers.NewEventualConsistencyPoller(10, func(pollerCtx context.Context) (*http.Response, error) {
		resp, err := client.Get(pollerCtx, id)
		return resp.HttpResponse, err
	}, &custompollers.EventualConsistencyPollerOptions{
		Interval:              15 * time.Second,
		TargetStatusCode:      pointer.To(http.StatusOK),
		RetryErrorStatusCodes: []int{http.StatusNotFound},
	})
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for Monitor %s to finish provisioning: %s", id, err)
	}

	d.SetId(id.ID())

	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceMonitorMetricAlertRead(d, meta)
}
