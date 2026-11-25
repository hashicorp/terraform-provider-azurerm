// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_pipeline

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/pipelines"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataFactoryPipelineCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.PipelinesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataFactoryId, err := factories.ParseFactoryID(d.Get("data_factory_id").(string))
	if err != nil {
		return err
	}

	id := pipelines.NewPipelineID(subscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id, pipelines.DefaultGetOperationOptions())
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_data_factory_pipeline", id.ID())
			}
		}
	}

	payload := pipelines.PipelineResource{
		Properties: pipelines.Pipeline{
			Description: pointer.To(d.Get("description").(string)),
			Parameters:  expandDataFactoryPipelineParameters(d.Get("parameters").(map[string]any)),
			Variables:   expandDataFactoryPipelineVariables(d.Get("variables").(map[string]any)),
		},
	}

	if v, ok := d.GetOk("activities_json"); ok && v.(string) != "" {
		activitiesJson, err := pipelines.UnmarshalActivityImplementation([]byte(fmt.Sprintf(`{ "activities": %s }`, v.(string))))
		if err != nil {
			return fmt.Errorf("unmarshaling `activities_json`: %+v", err)
		}
		rawActivities, ok := activitiesJson.(pipelines.RawActivityImpl)
		if !ok {
			return fmt.Errorf("expected `activities_json` to be of type `RawActivityImpl`")
		}

		activities := make([]pipelines.Activity, 0)
		acts, ok := rawActivities.Values["activities"]
		if !ok {
			return fmt.Errorf("`activities` was not found in the unmarshaled `activities_json`")
		}

		for _, activity := range acts.([]any) {
			act, err := json.Marshal(activity)
			if err != nil {
				return fmt.Errorf("marshaling activity %+v: %+v", activity, err)
			}
			a, err := pipelines.UnmarshalActivityImplementation(act)
			if err != nil {
				return fmt.Errorf("unmarshaling activity %+v: %+v", act, err)
			}
			activities = append(activities, a)
		}

		payload.Properties.Activities = pointer.To(activities)
	}

	annotations := make([]any, 0)
	if v, ok := d.GetOk("annotations"); ok {
		annotations = v.([]any)
	}
	payload.Properties.Annotations = &annotations

	if v, ok := d.GetOk("concurrency"); ok {
		payload.Properties.Concurrency = pointer.To(int64(v.(int)))
	}

	if v, ok := d.GetOk("monitor_metrics_after_duration"); ok {
		payload.Properties.Policy = &pipelines.PipelinePolicy{
			ElapsedTimeMetric: &pipelines.PipelineElapsedTimeMetricPolicy{
				Duration: pointer.To(v),
			},
		}
	}

	if v, ok := d.GetOk("folder"); ok {
		payload.Properties.Folder = &pipelines.PipelineFolder{
			Name: pointer.To(v.(string)),
		}
	}

	if _, err := client.CreateOrUpdate(ctx, id, payload, pipelines.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDataFactoryPipelineRead(d, meta)
}
